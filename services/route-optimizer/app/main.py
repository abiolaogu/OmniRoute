"""
FastAPI server for Route Optimization Service.
Provides REST API for VRP solving.
"""
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List, Optional
import uvicorn

from vrp_solver import (
    VRPSolver, VRPConfig, Location, Vehicle,
    calculate_distance_matrix, VRPSolution
)

app = FastAPI(
    title="OmniRoute Route Optimizer",
    description="Vehicle Routing Problem solver using Google OR-Tools",
    version="1.0.0"
)


class LocationRequest(BaseModel):
    id: str
    latitude: float
    longitude: float
    demand: int = 1
    time_window_start: Optional[int] = None
    time_window_end: Optional[int] = None
    service_time: int = 10


class VehicleRequest(BaseModel):
    id: str
    capacity: int


class OptimizeRequest(BaseModel):
    locations: List[LocationRequest]
    vehicles: List[VehicleRequest]
    time_limit_seconds: int = 30
    use_time_windows: bool = True


class RouteResponse(BaseModel):
    vehicle_id: str
    stop_ids: List[str]
    total_distance_meters: int
    total_time_minutes: int
    total_load: int


class OptimizeResponse(BaseModel):
    routes: List[RouteResponse]
    total_distance_meters: int
    total_time_minutes: int
    dropped_locations: List[str]
    computation_time_ms: int


@app.get("/health")
def health():
    return {"status": "healthy", "service": "route-optimizer"}


@app.get("/ready")
def ready():
    return {"status": "ready"}


@app.post("/optimize", response_model=OptimizeResponse)
def optimize_routes(request: OptimizeRequest):
    """Optimize delivery routes using VRP solver."""
    if len(request.locations) < 2:
        raise HTTPException(
            status_code=400,
            detail="At least 2 locations required (depot + 1 delivery)"
        )

    if len(request.vehicles) < 1:
        raise HTTPException(
            status_code=400,
            detail="At least 1 vehicle required"
        )

    # Convert to internal types
    locations = [
        Location(
            id=loc.id,
            latitude=loc.latitude,
            longitude=loc.longitude,
            demand=loc.demand,
            time_window_start=loc.time_window_start,
            time_window_end=loc.time_window_end,
            service_time=loc.service_time
        )
        for loc in request.locations
    ]

    vehicles = [
        Vehicle(
            id=v.id,
            capacity=v.capacity,
            start_location=locations[0]
        )
        for v in request.vehicles
    ]

    # Calculate distance matrix
    distance_matrix = calculate_distance_matrix(locations)

    # Configure and solve
    config = VRPConfig(
        time_limit_seconds=request.time_limit_seconds,
        use_time_windows=request.use_time_windows
    )
    solver = VRPSolver(config)
    solution = solver.solve(locations, vehicles, distance_matrix)

    # Convert response
    routes = [
        RouteResponse(
            vehicle_id=route.vehicle_id,
            stop_ids=[s.id for s in route.stops],
            total_distance_meters=route.total_distance,
            total_time_minutes=route.total_time,
            total_load=route.total_load
        )
        for route in solution.routes
    ]

    return OptimizeResponse(
        routes=routes,
        total_distance_meters=solution.total_distance,
        total_time_minutes=solution.total_time,
        dropped_locations=solution.dropped_locations,
        computation_time_ms=solution.computation_time_ms
    )


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8088)
