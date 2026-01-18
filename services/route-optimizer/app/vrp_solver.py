"""
OmniRoute Route Optimization Service
Uses Google OR-Tools for Vehicle Routing Problem (VRP) solving.
"""
from typing import List, Optional, Tuple
from dataclasses import dataclass
from ortools.constraint_solver import routing_enums_pb2
from ortools.constraint_solver import pywrapcp
import numpy as np


@dataclass
class Location:
    """Represents a delivery location."""
    id: str
    latitude: float
    longitude: float
    demand: int = 1
    time_window_start: Optional[int] = None  # Minutes from start of day
    time_window_end: Optional[int] = None
    service_time: int = 10  # Minutes to serve


@dataclass
class Vehicle:
    """Represents a delivery vehicle."""
    id: str
    capacity: int
    start_location: Location
    end_location: Optional[Location] = None
    max_distance: Optional[int] = None  # Meters
    max_time: Optional[int] = None  # Minutes


@dataclass
class VRPConfig:
    """Configuration for VRP solver."""
    time_limit_seconds: int = 30
    first_solution_strategy: str = "PATH_CHEAPEST_ARC"
    local_search_metaheuristic: str = "GUIDED_LOCAL_SEARCH"
    use_time_windows: bool = True
    use_capacity: bool = True


@dataclass
class Route:
    """Represents an optimized route."""
    vehicle_id: str
    stops: List[Location]
    total_distance: int  # Meters
    total_time: int  # Minutes
    total_load: int


@dataclass
class VRPSolution:
    """Solution to the VRP problem."""
    routes: List[Route]
    total_distance: int
    total_time: int
    dropped_locations: List[str]
    computation_time_ms: int


class VRPSolver:
    """Vehicle Routing Problem Solver using Google OR-Tools."""

    def __init__(self, config: VRPConfig):
        self.config = config

    def solve(
        self,
        locations: List[Location],
        vehicles: List[Vehicle],
        distance_matrix: List[List[int]],
        time_matrix: Optional[List[List[int]]] = None
    ) -> VRPSolution:
        """
        Solve the VRP problem.

        Args:
            locations: List of delivery locations (index 0 is depot)
            vehicles: List of available vehicles
            distance_matrix: Distance matrix in meters
            time_matrix: Optional time matrix in minutes

        Returns:
            VRPSolution with optimized routes
        """
        import time
        start_time = time.time()

        num_locations = len(locations)
        num_vehicles = len(vehicles)

        # Create routing index manager
        manager = pywrapcp.RoutingIndexManager(
            num_locations,
            num_vehicles,
            0  # Depot index
        )

        # Create routing model
        routing = pywrapcp.RoutingModel(manager)

        # Distance callback
        def distance_callback(from_index: int, to_index: int) -> int:
            from_node = manager.IndexToNode(from_index)
            to_node = manager.IndexToNode(to_index)
            return distance_matrix[from_node][to_node]

        transit_callback_index = routing.RegisterTransitCallback(distance_callback)
        routing.SetArcCostEvaluatorOfAllVehicles(transit_callback_index)

        # Add distance dimension
        dimension_name = "Distance"
        routing.AddDimension(
            transit_callback_index,
            0,  # No slack
            100000,  # Maximum distance per vehicle (100km)
            True,  # Start cumul to zero
            dimension_name
        )

        # Add capacity constraints
        if self.config.use_capacity:
            def demand_callback(from_index: int) -> int:
                from_node = manager.IndexToNode(from_index)
                return locations[from_node].demand

            demand_callback_index = routing.RegisterUnaryTransitCallback(demand_callback)
            
            vehicle_capacities = [v.capacity for v in vehicles]
            routing.AddDimensionWithVehicleCapacity(
                demand_callback_index,
                0,  # No slack
                vehicle_capacities,
                True,
                "Capacity"
            )

        # Add time window constraints
        if self.config.use_time_windows and time_matrix:
            def time_callback(from_index: int, to_index: int) -> int:
                from_node = manager.IndexToNode(from_index)
                to_node = manager.IndexToNode(to_index)
                return time_matrix[from_node][to_node] + locations[from_node].service_time

            time_callback_index = routing.RegisterTransitCallback(time_callback)

            routing.AddDimension(
                time_callback_index,
                30,  # Allow 30 min slack
                480,  # 8 hours maximum
                False,
                "Time"
            )

            time_dimension = routing.GetDimensionOrDie("Time")

            # Add time windows for each location
            for location_idx, location in enumerate(locations):
                if location.time_window_start is not None and location.time_window_end is not None:
                    index = manager.NodeToIndex(location_idx)
                    time_dimension.CumulVar(index).SetRange(
                        location.time_window_start,
                        location.time_window_end
                    )

        # Allow dropping locations with penalty
        penalty = 100000
        for i in range(1, num_locations):  # Skip depot
            routing.AddDisjunction([manager.NodeToIndex(i)], penalty)

        # Set search parameters
        search_parameters = pywrapcp.DefaultRoutingSearchParameters()
        search_parameters.first_solution_strategy = self._get_first_solution_strategy()
        search_parameters.local_search_metaheuristic = self._get_metaheuristic()
        search_parameters.time_limit.seconds = self.config.time_limit_seconds

        # Solve
        solution = routing.SolveWithParameters(search_parameters)

        computation_time = int((time.time() - start_time) * 1000)

        if solution:
            return self._extract_solution(
                manager, routing, solution, locations, vehicles, computation_time
            )
        else:
            return VRPSolution(
                routes=[],
                total_distance=0,
                total_time=0,
                dropped_locations=[loc.id for loc in locations[1:]],
                computation_time_ms=computation_time
            )

    def _extract_solution(
        self,
        manager: pywrapcp.RoutingIndexManager,
        routing: pywrapcp.RoutingModel,
        solution: pywrapcp.Assignment,
        locations: List[Location],
        vehicles: List[Vehicle],
        computation_time: int
    ) -> VRPSolution:
        """Extract solution from OR-Tools result."""
        routes = []
        total_distance = 0
        total_time = 0
        visited = set()

        for vehicle_idx in range(len(vehicles)):
            route_stops = []
            route_distance = 0
            route_load = 0

            index = routing.Start(vehicle_idx)
            while not routing.IsEnd(index):
                node_index = manager.IndexToNode(index)
                visited.add(node_index)
                
                if node_index > 0:  # Skip depot
                    route_stops.append(locations[node_index])
                    route_load += locations[node_index].demand

                previous_index = index
                index = solution.Value(routing.NextVar(index))
                route_distance += routing.GetArcCostForVehicle(
                    previous_index, index, vehicle_idx
                )

            if route_stops:  # Only add non-empty routes
                routes.append(Route(
                    vehicle_id=vehicles[vehicle_idx].id,
                    stops=route_stops,
                    total_distance=route_distance,
                    total_time=route_distance // 50,  # Rough estimate
                    total_load=route_load
                ))
                total_distance += route_distance

        # Find dropped locations
        dropped = [
            locations[i].id
            for i in range(1, len(locations))
            if i not in visited
        ]

        return VRPSolution(
            routes=routes,
            total_distance=total_distance,
            total_time=total_distance // 50,
            dropped_locations=dropped,
            computation_time_ms=computation_time
        )

    def _get_first_solution_strategy(self) -> int:
        """Get OR-Tools first solution strategy."""
        strategies = {
            "AUTOMATIC": routing_enums_pb2.FirstSolutionStrategy.AUTOMATIC,
            "PATH_CHEAPEST_ARC": routing_enums_pb2.FirstSolutionStrategy.PATH_CHEAPEST_ARC,
            "PATH_MOST_CONSTRAINED_ARC": routing_enums_pb2.FirstSolutionStrategy.PATH_MOST_CONSTRAINED_ARC,
            "SAVINGS": routing_enums_pb2.FirstSolutionStrategy.SAVINGS,
            "SWEEP": routing_enums_pb2.FirstSolutionStrategy.SWEEP,
            "CHRISTOFIDES": routing_enums_pb2.FirstSolutionStrategy.CHRISTOFIDES,
        }
        return strategies.get(
            self.config.first_solution_strategy,
            routing_enums_pb2.FirstSolutionStrategy.PATH_CHEAPEST_ARC
        )

    def _get_metaheuristic(self) -> int:
        """Get OR-Tools metaheuristic."""
        metaheuristics = {
            "AUTOMATIC": routing_enums_pb2.LocalSearchMetaheuristic.AUTOMATIC,
            "GREEDY_DESCENT": routing_enums_pb2.LocalSearchMetaheuristic.GREEDY_DESCENT,
            "GUIDED_LOCAL_SEARCH": routing_enums_pb2.LocalSearchMetaheuristic.GUIDED_LOCAL_SEARCH,
            "SIMULATED_ANNEALING": routing_enums_pb2.LocalSearchMetaheuristic.SIMULATED_ANNEALING,
            "TABU_SEARCH": routing_enums_pb2.LocalSearchMetaheuristic.TABU_SEARCH,
        }
        return metaheuristics.get(
            self.config.local_search_metaheuristic,
            routing_enums_pb2.LocalSearchMetaheuristic.GUIDED_LOCAL_SEARCH
        )


def calculate_distance_matrix(locations: List[Location]) -> List[List[int]]:
    """
    Calculate Haversine distance matrix between all locations.
    Returns distance in meters.
    """
    n = len(locations)
    matrix = [[0] * n for _ in range(n)]

    for i in range(n):
        for j in range(n):
            if i != j:
                matrix[i][j] = haversine_distance(
                    locations[i].latitude, locations[i].longitude,
                    locations[j].latitude, locations[j].longitude
                )

    return matrix


def haversine_distance(lat1: float, lon1: float, lat2: float, lon2: float) -> int:
    """Calculate Haversine distance in meters."""
    R = 6371000  # Earth radius in meters

    phi1 = np.radians(lat1)
    phi2 = np.radians(lat2)
    delta_phi = np.radians(lat2 - lat1)
    delta_lambda = np.radians(lon2 - lon1)

    a = np.sin(delta_phi/2)**2 + np.cos(phi1) * np.cos(phi2) * np.sin(delta_lambda/2)**2
    c = 2 * np.arctan2(np.sqrt(a), np.sqrt(1-a))

    return int(R * c)


# Example usage
if __name__ == "__main__":
    # Create locations (depot + deliveries)
    locations = [
        Location(id="depot", latitude=6.5244, longitude=3.3792),  # Lagos
        Location(id="loc1", latitude=6.4474, longitude=3.3903, demand=2, time_window_start=60, time_window_end=180),
        Location(id="loc2", latitude=6.4698, longitude=3.5852, demand=1, time_window_start=120, time_window_end=240),
        Location(id="loc3", latitude=6.5355, longitude=3.3087, demand=3, time_window_start=60, time_window_end=300),
        Location(id="loc4", latitude=6.4281, longitude=3.4219, demand=2),
        Location(id="loc5", latitude=6.5915, longitude=3.3449, demand=1),
    ]

    # Create vehicles
    vehicles = [
        Vehicle(id="v1", capacity=5, start_location=locations[0]),
        Vehicle(id="v2", capacity=5, start_location=locations[0]),
    ]

    # Calculate distance matrix
    distance_matrix = calculate_distance_matrix(locations)

    # Solve
    config = VRPConfig(time_limit_seconds=10, use_time_windows=False)
    solver = VRPSolver(config)
    solution = solver.solve(locations, vehicles, distance_matrix)

    print(f"Total distance: {solution.total_distance / 1000:.2f} km")
    print(f"Computation time: {solution.computation_time_ms} ms")
    for route in solution.routes:
        print(f"Vehicle {route.vehicle_id}: {[s.id for s in route.stops]}")
