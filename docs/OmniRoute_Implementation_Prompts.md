# OmniRoute B2B FMCG Platform - Complete Implementation Prompts

## Quick Reference Guide

| Layer | Prompts | Description |
|-------|---------|-------------|
| L1: Commerce Core | 4 | Catalog, Orders, Inventory, Customers |
| L2: Gig Platform | 4 | Workers, Tasks, Earnings, Careers |
| L3: Distribution | 4 | Routes, WMS, Restocking, Fleet |
| L4: Accessibility | 4 | Voice, USSD, NLU, WhatsApp |
| L5: Social Commerce | 3 | Groups, Reputation, Referrals |
| L6: Intelligence | 3 | Market Intel, Analytics, AI Insights |
| L7: Finance | 4 | Lending, Payments, Wallets, Credit Scoring |
| Mobile Apps | 3 | Customer, Worker, Admin |
| Security | 2 | Auth, Privacy |
| DevOps/MLOps | 3 | Terraform, Helm, ML Pipelines |
| Observability | 1 | Full Stack |
| Testing/Docs | 2 | Framework, Documentation |

**Total: 37 Production-Ready Prompts**

---

# LAYER 1: COMMERCE CORE

## L1-P01: Product Catalog Service (Go)

```
PROJECT: OmniRoute - Product Catalog Service

CONTEXT:
Build a high-performance Product Catalog service for FMCG distribution supporting:
- Multi-tenant product hierarchies (Category > Subcategory > Brand > SKU)
- Dynamic pricing rules with regional variations
- Unit of Measure conversions (case, pack, unit)
- Product bundles and promotions
- Real-time inventory visibility across warehouses
- Product images and rich media management
- Barcode/QR code integration

TECHNICAL REQUIREMENTS:
- Language: Go 1.22+
- Database: PostgreSQL 16 with JSONB for flexible attributes
- Cache: Redis Cluster for hot product data
- Search: Elasticsearch 8.x for full-text product search
- Storage: MinIO for product images
- Protocol: gRPC for internal, REST for external

DOMAIN MODEL:
- Product aggregate with SKU as identity
- PriceList value object with validity periods
- InventoryPosition entity per warehouse
- ProductBundle aggregate for composite products

PERFORMANCE TARGETS:
- Catalog queries: <10ms p99
- Search results: <50ms p99
- Support 1M+ SKUs per tenant
- 10K concurrent users per region

CODE STRUCTURE:
cmd/
  catalog/
    main.go
internal/
  domain/
    product.go
    price.go
    inventory.go
  repository/
    product_repo.go
    price_repo.go
  handler/
    grpc_handler.go
    http_handler.go
  search/
    elasticsearch.go
  cache/
    redis_cache.go
pkg/
  proto/
    catalog.proto

TEST CASES:
1. Product CRUD operations with validation
2. Multi-tenant isolation verification
3. Price calculation with tier discounts
4. Search relevance scoring
5. Concurrent inventory updates
6. Cache invalidation on price changes
7. Bulk import performance (1000 SKUs/second)
8. Image upload and CDN integration

EXPECTED DELIVERABLES:
1. Complete Go project structure
2. Domain models with validation
3. Repository implementations
4. gRPC service definitions (.proto files)
5. REST API handlers with OpenAPI spec
6. Elasticsearch indexing and search
7. Redis caching layer
8. Unit tests (>85% coverage)
9. Integration tests with testcontainers
10. Kubernetes deployment manifests
11. OpenTelemetry instrumentation
12. Database migrations

Please implement following DDD principles, XP practices (TDD), and include comprehensive error handling.
```

---

## L1-P02: Order Management Service (Go)

```
PROJECT: OmniRoute - Order Management Service

CONTEXT:
Build a comprehensive Order Management System handling:
- Multi-channel order capture (Web, Mobile, USSD, Voice, WhatsApp)
- B2B order workflows with approval chains
- Split shipments from multiple warehouses
- Order modifications and cancellations
- Returns and refunds processing
- Real-time order status tracking
- Integration with Temporal for saga orchestration

TECHNICAL REQUIREMENTS:
- Language: Go 1.22+
- Workflow: Temporal for order saga orchestration
- Event Sourcing: EventStoreDB for order history
- Database: PostgreSQL for projections
- Messaging: Apache Kafka for event streaming

ORDER STATES:
Draft -> Submitted -> Approved -> Confirmed -> Picking -> Packed -> Shipped -> Delivered -> Completed

SAGA WORKFLOW:
OrderSagaWorkflow:
  1. ValidateOrder - Check product availability, customer eligibility
  2. ReserveInventory - Hold inventory across warehouses
  3. CalculatePricing - Apply discounts, taxes, shipping
  4. ProcessPayment - Charge customer or create invoice
  5. CreateShipment - Generate shipping labels, assign carrier
  6. NotifyParties - Send confirmations to all stakeholders

COMPENSATING ACTIONS:
- ReleaseInventory - On failure after reservation
- RefundPayment - On failure after payment
- CancelShipment - On failure after shipment creation

DOMAIN MODEL:
type Order struct {
    ID              uuid.UUID
    TenantID        uuid.UUID
    CustomerID      uuid.UUID
    Channel         OrderChannel
    Status          OrderStatus
    Lines           []OrderLine
    ShippingAddress Address
    BillingAddress  Address
    TotalAmount     Money
    Discounts       []Discount
    Payments        []Payment
    Shipments       []Shipment
    ApprovalChain   []Approval
    Events          []DomainEvent
    Version         int
}

TEST CASES:
1. Happy path order completion through all states
2. Inventory reservation failure with compensation
3. Payment failure with rollback
4. Partial shipment handling
5. Order modification mid-saga
6. Concurrent order conflicts (optimistic locking)
7. Idempotency on retry
8. Event replay and projection rebuild

EXPECTED DELIVERABLES:
1. Temporal workflow definitions
2. Activity implementations with retry policies
3. Event-sourced Order aggregate
4. Read model projections
5. gRPC and REST APIs
6. Kafka event publishers
7. Unit and integration tests
8. Deployment manifests
```

---

## L1-P03: Inventory Management Service (Go)

```
PROJECT: OmniRoute - Inventory Management Service

CONTEXT:
Build a distributed Inventory Management service supporting:
- Real-time inventory across multiple warehouses
- Reservation and allocation strategies
- Batch and expiry tracking (FEFO/FIFO)
- Inventory adjustments and cycle counts
- Warehouse zone management
- Reorder point calculations
- ATP (Available-to-Promise) engine

TECHNICAL REQUIREMENTS:
- Language: Go 1.22+
- Database: PostgreSQL with row-level locking for consistency
- Cache: Redis for ATP calculations
- Events: Kafka for inventory changes
- Transactions: Distributed transactions with saga pattern

INVENTORY TYPES:
- On-Hand: Physical inventory in warehouse
- Reserved: Allocated to confirmed orders
- Available: On-Hand minus Reserved
- In-Transit: Being transferred between locations
- Quarantine: Quality hold, not available for sale

ALLOCATION STRATEGIES:
- FIFO: First In First Out (by receipt date)
- FEFO: First Expiry First Out (by expiry date)
- Proximity: Nearest warehouse to customer
- Cost: Lowest cost warehouse (including shipping)

DOMAIN MODEL:
type InventoryPosition struct {
    ID           uuid.UUID
    TenantID     uuid.UUID
    WarehouseID  uuid.UUID
    ProductID    uuid.UUID
    LocationID   uuid.UUID
    BatchNumber  string
    ExpiryDate   *time.Time
    OnHand       decimal.Decimal
    Reserved     decimal.Decimal
    InTransit    decimal.Decimal
    Quarantine   decimal.Decimal
    UnitCost     Money
    LastCountAt  time.Time
}

ATP ENGINE:
type ATPEngine interface {
    CalculateATP(ctx context.Context, productID uuid.UUID, 
        warehouseIDs []uuid.UUID, date time.Time) (decimal.Decimal, error)
    ProjectATP(ctx context.Context, productID uuid.UUID, 
        days int) ([]ATPProjection, error)
}

// ATP = On-Hand - Reserved + Incoming (POs) - Outgoing (SOs)

TEST CASES:
1. Concurrent reservation conflicts (pessimistic/optimistic locking)
2. FEFO allocation accuracy with multiple batches
3. Cross-warehouse ATP calculation
4. Batch expiry enforcement
5. Cycle count reconciliation
6. Real-time inventory streaming to Kafka
7. Performance under high contention (1000 TPS)

EXPECTED DELIVERABLES:
1. Domain models and repositories
2. ATP calculation engine with caching
3. Allocation service with strategy pattern
4. Event-driven updates via Kafka
5. gRPC service definitions
6. Warehouse zone management APIs
7. Comprehensive tests
8. Performance benchmarks
```

---

## L1-P04: Customer Management Service (Go)

```
PROJECT: OmniRoute - Customer Management Service

CONTEXT:
Build a B2B Customer Management service handling:
- Multi-level customer hierarchies (Distributor > Retailer > Outlet)
- Customer segmentation and tiers
- Credit limits and payment terms
- Sales territory assignments
- Customer onboarding workflows
- KYC verification integration
- Customer 360 view aggregation

TECHNICAL REQUIREMENTS:
- Language: Go 1.22+
- Database: PostgreSQL with hierarchical queries (recursive CTEs)
- Search: Elasticsearch for customer lookup
- Events: Kafka for customer events
- Integration: External KYC providers

CUSTOMER TIERS:
- Platinum: >$100K/month volume, 60-day payment terms, 5% discount
- Gold: >$50K/month volume, 45-day payment terms, 3% discount
- Silver: >$10K/month volume, 30-day payment terms, 1% discount
- Bronze: Standard terms, 15-day payment, no discount

KYC LEVELS:
- Basic: Phone verification + Government ID
- Standard: + Business registration documents
- Enhanced: + Bank account verification + Trade references

DOMAIN MODEL:
type Customer struct {
    ID             uuid.UUID
    TenantID       uuid.UUID
    ParentID       *uuid.UUID  // For hierarchy
    Type           CustomerType // Distributor, Retailer, Outlet
    Tier           CustomerTier
    BusinessName   string
    TaxID          string
    ContactPerson  ContactInfo
    Addresses      []Address
    CreditLimit    Money
    CreditUsed     Money
    PaymentTerms   int // Days
    KYCStatus      KYCStatus
    KYCLevel       KYCLevel
    SalesRepID     uuid.UUID
    TerritoryID    uuid.UUID
    Metadata       map[string]interface{}
}

HIERARCHY QUERIES:
- Get all descendants of a distributor
- Get the parent chain for an outlet
- Calculate aggregate volume for hierarchy
- Enforce credit limits across hierarchy

TEST CASES:
1. Customer hierarchy traversal with recursive queries
2. Credit limit enforcement at order time
3. KYC workflow completion with external verification
4. Territory assignment rules and validation
5. Tier upgrade/downgrade based on volume
6. Customer merge handling (duplicate resolution)
7. GDPR data export (right to portability)

EXPECTED DELIVERABLES:
1. Customer domain models with hierarchy support
2. Hierarchy management service
3. KYC integration service
4. Credit management service
5. Territory assignment service
6. gRPC/REST APIs
7. Kafka event publishers
8. Tests and documentation
```

---

# LAYER 2: GIG PLATFORM

## L2-P01: Worker Management Service (Go)

```
PROJECT: OmniRoute - Worker Management Service

CONTEXT:
Build a comprehensive Gig Worker Management platform supporting:
- Multi-role workers (Driver, Picker, Merchandiser, Sales Rep)
- Skill-based matching and certification
- Worker onboarding and verification
- Performance scoring and ratings
- Availability and shift management
- Worker wallet and earnings
- Benefits and insurance integration
- Career progression pathways

TECHNICAL REQUIREMENTS:
- Language: Go 1.22+
- Database: PostgreSQL
- Geospatial: PostGIS for location services
- Real-time: Redis Pub/Sub for availability updates
- Messaging: Kafka for worker events

WORKER ROLES:
- Delivery Driver: Last-mile delivery, requires license
- Van Sales Rep: Route sales execution, requires sales training
- Picker: Warehouse order picking, requires WMS certification
- Merchandiser: In-store execution, requires brand training
- Territory Sales Rep: B2B sales, requires sales certification
- Multi-role: Can perform multiple roles based on certifications

VERIFICATION REQUIREMENTS:
- Identity: Government ID + Selfie verification
- Driver's License: For driver roles (validated against registry)
- Background Check: Criminal record check
- Vehicle Inspection: For drivers with own vehicles
- Skills Assessment: Role-specific competency tests

DOMAIN MODEL:
type Worker struct {
    ID              uuid.UUID
    TenantID        uuid.UUID
    Profile         WorkerProfile
    Roles           []WorkerRole
    Skills          []Skill
    Certifications  []Certification
    Vehicle         *Vehicle
    WorkingAreas    []GeographicArea  // PostGIS polygons
    Availability    []AvailabilitySlot
    PerformanceScore float64
    Status          WorkerStatus
    Wallet          Wallet
    Documents       []Document
    OnboardingStatus OnboardingStatus
}

type WorkerRole struct {
    Role           RoleType
    Status         RoleStatus  // Pending, Active, Suspended
    CertifiedAt    *time.Time
    ExpiresAt      *time.Time
    Restrictions   []string
    PerformanceByRole float64
}

MATCHING ALGORITHM:
func FindBestWorkers(task Task, count int) []WorkerMatch {
    // 1. Filter by required role and skills
    // 2. Filter by geographic coverage (PostGIS intersection)
    // 3. Filter by current availability
    // 4. Score by: performance (40%), proximity (30%), 
    //              workload balance (20%), cost (10%)
    // 5. Return top N matches with scores
}

TEST CASES:
1. Multi-role worker registration and role activation
2. Skill certification workflow with expiry
3. Availability conflict detection
4. Location-based matching with PostGIS
5. Performance score calculation from completed tasks
6. Wallet transaction integrity
7. Concurrent task assignment handling
8. Role transition workflows

EXPECTED DELIVERABLES:
1. Worker domain models with role management
2. Matching algorithm service
3. Availability management service
4. Performance scoring engine
5. Wallet service integration
6. Document verification service
7. gRPC/REST APIs
8. Real-time notification service
9. Tests and benchmarks
```

---

## L2-P02: Task Assignment Service (Go)

```
PROJECT: OmniRoute - Task Assignment Service

CONTEXT:
Build an intelligent Task Assignment engine handling:
- Dynamic task creation from orders/routes
- Real-time worker matching
- Batch assignment optimization
- Task handoffs and reassignment
- SLA tracking and escalation
- Task completion verification
- Photo/signature capture
- Proof of delivery

TECHNICAL REQUIREMENTS:
- Language: Go 1.22+
- Optimization: Google OR-Tools for batch assignment
- Real-time: WebSocket for live updates
- Storage: MinIO for proof-of-delivery images
- Geospatial: PostGIS for proximity calculations

TASK TYPES:
- Delivery: Last-mile order delivery
- Pickup: Collection from supplier or returns
- VanSale: Route-based selling (multiple stops)
- Merchandising: In-store product placement
- Survey: Market research data collection
- Inspection: Quality audit tasks

ASSIGNMENT STRATEGIES:
- Nearest Available: Minimize travel distance
- Lowest Cost: Consider worker rates and travel
- Best Performance: Prioritize high-rated workers
- Load Balancing: Distribute tasks evenly
- Skills Required: Match specialized skills

BATCH OPTIMIZATION (OR-Tools):
// Solve assignment problem for multiple tasks/workers
func OptimizeAssignments(tasks []Task, workers []Worker) []Assignment {
    // Objective: Minimize total cost (time + distance + worker rate)
    // Constraints:
    //   - Worker capacity (max tasks per shift)
    //   - Time windows (pickup/delivery constraints)
    //   - Skill requirements
    //   - Geographic zones
    // Returns: Optimal task-to-worker assignments
}

DOMAIN MODEL:
type Task struct {
    ID              uuid.UUID
    TenantID        uuid.UUID
    Type            TaskType
    Status          TaskStatus
    Priority        int  // 1-5, affects SLA
    AssignedWorkerID *uuid.UUID
    Location        GeoPoint
    ScheduledTime   TimeWindow
    ActualStartTime *time.Time
    ActualEndTime   *time.Time
    Instructions    string
    Proofs          []ProofOfCompletion
    SLA             SLAConfig
    Compensation    Compensation
    Dependencies    []uuid.UUID  // Tasks that must complete first
}

type ProofOfCompletion struct {
    Type       ProofType  // Photo, Signature, Barcode, GPS
    Data       []byte
    CapturedAt time.Time
    Location   GeoPoint
    Verified   bool
}

SLA CONFIGURATION:
type SLAConfig struct {
    ResponseTime   time.Duration  // Time to accept
    CompletionTime time.Duration  // Time to complete
    EscalationRules []EscalationRule
}

TEST CASES:
1. Single task assignment with matching
2. Batch optimization accuracy (vs brute force)
3. SLA breach escalation workflow
4. Task reassignment on worker unavailability
5. Proof-of-delivery validation
6. Concurrent assignment conflicts
7. Worker offline handling
8. Performance under load (10K tasks/hour)

EXPECTED DELIVERABLES:
1. Task domain models
2. Assignment optimizer with OR-Tools
3. Real-time matching service
4. SLA monitoring service
5. Proof collection and verification
6. WebSocket service for live updates
7. APIs and events
8. Comprehensive tests
```

---

## L2-P03: Earnings & Benefits Service (Go)

```
PROJECT: OmniRoute - Earnings & Benefits Service

CONTEXT:
Build a Worker Earnings and Benefits platform supporting:
- Flexible compensation models (per-task, hourly, commission)
- Real-time earnings calculation
- Instant payout to mobile money/bank
- Bonus and incentive programs
- Tax withholding calculations
- Insurance and benefits enrollment
- Savings programs
- Loan repayment deductions

TECHNICAL REQUIREMENTS:
- Language: Go 1.22+
- Database: PostgreSQL with full audit trail
- Payments: Mobile money APIs (M-Pesa, MTN MoMo, Airtel Money)
- Compliance: Tax calculation engine per jurisdiction

COMPENSATION MODELS:
- Per-Task: Fixed rate per task type and complexity
- Hourly: Time-based with minimum guarantees
- Commission: Percentage of sales amount
- Hybrid: Base pay + performance bonus
- Surge: Dynamic multipliers based on demand/time

BENEFITS PACKAGES:
- Basic: Accident insurance only
- Standard: + Health coverage + Paid time off
- Premium: + Savings match + Training allowance + Equipment

DOMAIN MODEL:
type Earnings struct {
    ID            uuid.UUID
    WorkerID      uuid.UUID
    Period        EarningsPeriod
    TaskEarnings  []TaskEarning
    Bonuses       []Bonus
    Deductions    []Deduction
    GrossAmount   Money
    TaxWithheld   Money
    NetAmount     Money
    Status        EarningsStatus
    PayoutMethod  PaymentMethod
    PaidAt        *time.Time
}

type TaskEarning struct {
    TaskID          uuid.UUID
    TaskType        TaskType
    BaseAmount      Money
    SurgeMultiplier float64
    Tips            Money
    Commission      Money
    TotalAmount     Money
    CompletedAt     time.Time
}

PAYOUT FLOW:
func ProcessPayout(workerID uuid.UUID, method PaymentMethod) error {
    // 1. Calculate net earnings for period
    // 2. Apply tax withholdings
    // 3. Deduct benefits premiums
    // 4. Deduct loan repayments
    // 5. Validate minimum payout amount
    // 6. Initiate transfer via payment gateway
    // 7. Record transaction with reconciliation ID
    // 8. Notify worker of successful payout
}

TEST CASES:
1. Multi-model earnings calculation
2. Surge pricing accuracy during peak times
3. Tax withholding compliance per jurisdiction
4. Instant payout processing (<30 seconds)
5. Benefits enrollment and deduction
6. Loan deduction accuracy
7. Period closing and reconciliation
8. Concurrent payout handling

EXPECTED DELIVERABLES:
1. Earnings domain models
2. Calculation engine with all models
3. Payout service with multiple gateways
4. Tax calculator (pluggable per country)
5. Benefits management service
6. Mobile money integrations
7. APIs and reports
8. Complete audit trail
9. Comprehensive tests
```

---

## L2-P04: Career Progression Service (Go)

```
PROJECT: OmniRoute - Career Progression Service

CONTEXT:
Build a Worker Career Progression system enabling:
- Skill development pathways
- Certification and training management
- Performance-based advancement
- Role transition support
- Mentorship matching
- Achievement and badge system
- Entrepreneurship pathway (own delivery business/franchise)

TECHNICAL REQUIREMENTS:
- Language: Go 1.22+
- Database: PostgreSQL
- LMS Integration: For training content delivery
- Gamification: Achievement and badge engine

CAREER PATHS:
- Driver Path: Driver -> Senior Driver -> Team Lead -> Fleet Manager
- Warehouse Path: Picker -> Lead Picker -> Shift Supervisor -> Warehouse Manager
- Sales Path: Sales Rep -> Senior Rep -> Territory Manager -> Regional Manager
- Entrepreneur Path: Any Role -> Business Owner (Franchise)

ADVANCEMENT CRITERIA:
type AdvancementCriteria struct {
    MinimumTenure      time.Duration
    MinimumPerformance float64  // Score threshold
    RequiredCerts      []CertificationID
    RequiredTraining   []TrainingID
    PeerEndorsements   int
    IncidentFreeMonths int
    SpecialConditions  []Condition
}

DOMAIN MODEL:
type WorkerCareer struct {
    ID              uuid.UUID
    WorkerID        uuid.UUID
    CurrentPathID   uuid.UUID
    CurrentLevelID  uuid.UUID
    Progress        []LevelProgress
    Achievements    []Achievement
    Certifications  []Certification
    TrainingHistory []TrainingRecord
    MentorID        *uuid.UUID
    MenteeIDs       []uuid.UUID
    Points          int  // Gamification points
    Badges          []Badge
}

type Achievement struct {
    ID          uuid.UUID
    Type        AchievementType
    Name        string
    Description string
    Icon        string
    Points      int
    UnlockedAt  time.Time
    Criteria    string  // JSON criteria that was met
}

GAMIFICATION ENGINE:
- Daily challenges with rewards
- Weekly leaderboards
- Milestone achievements
- Streak bonuses
- Community recognition

TEST CASES:
1. Path enrollment and initial level assignment
2. Advancement criteria evaluation
3. Certification tracking with expiry alerts
4. Achievement unlocking triggers
5. Mentorship matching algorithm
6. Training completion synchronization
7. Multi-path transitions
8. Franchise qualification assessment

EXPECTED DELIVERABLES:
1. Career domain models
2. Progression evaluation engine
3. Achievement and badge system
4. Training management integration
5. Mentorship matching service
6. Gamification engine
7. APIs and worker portal
8. Comprehensive tests
```

---

# LAYER 3: DISTRIBUTION

## L3-P01: Route Optimization Service (Go + Python)

```
PROJECT: OmniRoute - Route Optimization Service

CONTEXT:
Build an AI-powered Route Optimization engine for:
- Vehicle Routing Problem (VRP) solving
- Dynamic route recalculation
- Multi-stop delivery optimization
- Van sales route planning
- Territory coverage optimization
- Traffic and time-window aware routing
- Capacity and weight constraints
- Return-to-depot optimization

TECHNICAL REQUIREMENTS:
- Orchestration: Go 1.22+ for API and coordination
- Optimization: Python with Google OR-Tools / VROOM
- Maps: Google Maps API / OpenRouteService
- Real-time: Redis for route caching
- Streaming: Kafka for route updates

VRP CONSTRAINTS:
- Vehicle capacity (weight, volume, item count)
- Time windows per delivery stop
- Driver hours of service regulations
- Vehicle type restrictions (refrigerated, size)
- Priority deliveries (express, same-day)
- Pickup-and-delivery pairs

OPTIMIZATION OBJECTIVES:
- Minimize total distance traveled
- Minimize total time (including traffic)
- Maximize on-time delivery rate
- Balance workload across drivers
- Minimize vehicle operating costs

PYTHON SOLVER (OR-Tools):
```python
from ortools.constraint_solver import routing_enums_pb2
from ortools.constraint_solver import pywrapcp

class VRPSolver:
    def __init__(self, config: VRPConfig):
        self.config = config
        
    def solve(self, problem: VRPProblem) -> VRPSolution:
        # Create routing index manager
        manager = pywrapcp.RoutingIndexManager(
            len(problem.locations),
            problem.num_vehicles,
            problem.depot_indices
        )
        
        # Create routing model
        routing = pywrapcp.RoutingModel(manager)
        
        # Distance callback
        def distance_callback(from_index, to_index):
            from_node = manager.IndexToNode(from_index)
            to_node = manager.IndexToNode(to_index)
            return problem.distance_matrix[from_node][to_node]
        
        transit_callback_index = routing.RegisterTransitCallback(distance_callback)
        routing.SetArcCostEvaluatorOfAllVehicles(transit_callback_index)
        
        # Add capacity constraint
        def demand_callback(from_index):
            from_node = manager.IndexToNode(from_index)
            return problem.demands[from_node]
        
        demand_callback_index = routing.RegisterUnaryTransitCallback(demand_callback)
        routing.AddDimensionWithVehicleCapacity(
            demand_callback_index, 0, problem.vehicle_capacities, 
            True, 'Capacity'
        )
        
        # Add time window constraints
        # ... (time dimension setup)
        
        # Solve
        search_parameters = pywrapcp.DefaultRoutingSearchParameters()
        search_parameters.first_solution_strategy = (
            routing_enums_pb2.FirstSolutionStrategy.PATH_CHEAPEST_ARC
        )
        search_parameters.local_search_metaheuristic = (
            routing_enums_pb2.LocalSearchMetaheuristic.GUIDED_LOCAL_SEARCH
        )
        search_parameters.time_limit.seconds = self.config.time_limit_seconds
        
        solution = routing.SolveWithParameters(search_parameters)
        return self._extract_solution(manager, routing, solution)
```

GO ORCHESTRATION SERVICE:
```go
type RouteOptimizerService struct {
    pythonClient  *PythonSolverClient
    mapsClient    MapsClient
    trafficClient TrafficClient
    cache         *redis.Client
}

func (s *RouteOptimizerService) OptimizeRoutes(ctx context.Context,
    req *OptimizeRequest) (*OptimizeResponse, error) {
    // 1. Fetch current traffic data
    // 2. Build distance/time matrix with traffic
    // 3. Call Python VRP solver via gRPC
    // 4. Enrich with turn-by-turn directions
    // 5. Cache routes for quick retrieval
    // 6. Return optimized routes
}
```

TEST CASES:
1. Simple VRP solving accuracy (vs known optimal)
2. Time window constraint satisfaction
3. Capacity constraint enforcement
4. Dynamic re-optimization on changes
5. Large problem scaling (1000+ stops)
6. Traffic-aware routing accuracy
7. Multi-depot scenarios
8. Performance benchmarks (solve time)

EXPECTED DELIVERABLES:
1. Python VRP solver with OR-Tools
2. Go orchestration service
3. Distance/time matrix builder
4. Traffic integration service
5. gRPC service definitions
6. Route visualization API
7. Redis caching layer
8. Performance benchmarks
9. Tests and documentation
```

---

## L3-P02: Warehouse Management Service (Go)

```
PROJECT: OmniRoute - Warehouse Management Service

CONTEXT:
Build a Warehouse Management System (WMS) supporting:
- Multi-warehouse operations
- Zone and location management
- Receiving and put-away workflows
- Pick path optimization
- Wave planning and execution
- Cross-docking support
- Returns processing
- Cycle counting
- Integration with automation (conveyors, pick-to-light)

TECHNICAL REQUIREMENTS:
- Language: Go 1.22+
- Database: PostgreSQL
- Real-time: Redis Streams for warehouse events
- Hardware: Integration with barcode scanners, RF devices
- Printing: Label generation for picking/shipping

WAREHOUSE ZONES:
- Receiving: Inbound staging and quality check
- Bulk Storage: High-density rack storage
- Forward Pick: Fast-moving items, easy access
- Cross-dock: Pass-through without storage
- Returns: RMA processing area
- Shipping: Outbound staging and loading

PICKING STRATEGIES:
- Discrete: One order at a time (simple, accurate)
- Batch: Multiple orders together (efficient)
- Wave: Scheduled batches with cutoff times
- Zone: Workers assigned to zones, pass orders
- Cluster: Group orders by location similarity

LOCATION MANAGEMENT:
type Location struct {
    ID              uuid.UUID
    WarehouseID     uuid.UUID
    ZoneID          uuid.UUID
    Aisle           string
    Rack            string
    Level           string
    Position        string
    Barcode         string
    Type            LocationType  // Bulk, Pick, Staging
    Capacity        LocationCapacity
    CurrentContents []LocationContent
    PickSequence    int  // For path optimization
}

PICK PATH OPTIMIZATION:
func OptimizePickPath(picks []PickTask) []PickTask {
    // Build warehouse graph from location topology
    // Apply TSP solver for shortest path
    // Consider zone transitions (minimize zone changes)
    // Return picks in optimized sequence
}

WAVE PLANNING:
type Wave struct {
    ID           uuid.UUID
    WarehouseID  uuid.UUID
    CutoffTime   time.Time
    ShipTime     time.Time
    Orders       []uuid.UUID
    Status       WaveStatus
    PickTasks    []PickTask
    PackTasks    []PackTask
    WorkerAssignments []WorkerAssignment
}

TEST CASES:
1. Complete receiving workflow with put-away
2. Put-away optimization (minimize travel)
3. Pick path efficiency measurement
4. Wave execution with deadlines
5. Cross-dock timing accuracy
6. Cycle count reconciliation
7. Returns processing workflow
8. Concurrent operations in same zone

EXPECTED DELIVERABLES:
1. WMS domain models
2. Location and zone management
3. Pick path optimization
4. Wave planning engine
5. RF device integration
6. Real-time dashboards
7. APIs and events
8. Tests and benchmarks
```

---

## L3-P03: Predictive Restocking Service (Python)

```
PROJECT: OmniRoute - Predictive Restocking Service

CONTEXT:
Build an AI-powered Predictive Restocking system:
- Demand forecasting per SKU/location
- Optimal reorder point calculation
- Safety stock optimization
- Lead time prediction
- Seasonal pattern detection
- Promotion impact modeling
- Automatic purchase order generation
- Supplier performance integration

TECHNICAL REQUIREMENTS:
- Language: Python 3.11+
- ML Framework: PyTorch + Prophet + XGBoost
- Feature Store: Feast
- Orchestration: Temporal for batch jobs
- Database: PostgreSQL + TimescaleDB for time-series

FORECASTING MODELS:
- Prophet: Captures seasonality and holidays
- LSTM/DeepAR: Sequential pattern learning
- XGBoost: Feature-rich gradient boosting
- Ensemble: Weighted combination for robustness

SAFETY STOCK FACTORS:
- Demand variability (standard deviation)
- Lead time variability
- Service level target (e.g., 95% fill rate)
- Supplier reliability score
- Storage costs (holding cost tradeoff)

IMPLEMENTATION:
```python
class DemandForecaster:
    def __init__(self, config: ForecastConfig):
        self.config = config
        self.feature_store = FeatureStore("feature_repo/")
        self.models = {}
        
    def train(self, product_id: str, history: pd.DataFrame):
        # Feature engineering
        features = self._extract_features(history)
        
        # Train Prophet for seasonality
        prophet_model = Prophet(
            yearly_seasonality=True,
            weekly_seasonality=True,
            daily_seasonality=False,
            holidays=self._get_holidays()
        )
        prophet_model.add_regressor('promotion')
        prophet_model.add_regressor('price')
        prophet_model.fit(history[['ds', 'y', 'promotion', 'price']])
        
        # Train XGBoost for external factors
        xgb_model = xgb.XGBRegressor(
            objective='reg:squarederror',
            n_estimators=100,
            max_depth=6,
            learning_rate=0.1
        )
        xgb_model.fit(features, history['y'])
        
        # Train LSTM for sequences
        lstm_model = self._train_lstm(history)
        
        self.models[product_id] = {
            'prophet': prophet_model,
            'xgboost': xgb_model,
            'lstm': lstm_model,
            'weights': [0.4, 0.35, 0.25]  # Ensemble weights
        }
        
    def forecast(self, product_id: str, horizon_days: int) -> Forecast:
        models = self.models[product_id]
        
        # Generate forecasts from each model
        prophet_forecast = models['prophet'].predict(future_df)
        xgb_forecast = models['xgboost'].predict(future_features)
        lstm_forecast = models['lstm'].predict(future_sequences)
        
        # Ensemble combination
        ensemble = (
            models['weights'][0] * prophet_forecast['yhat'] +
            models['weights'][1] * xgb_forecast +
            models['weights'][2] * lstm_forecast
        )
        
        # Calculate prediction intervals
        lower, upper = self._calculate_intervals(ensemble)
        
        return Forecast(
            product_id=product_id,
            predictions=ensemble,
            lower_bound=lower,
            upper_bound=upper,
            confidence=0.95
        )


class ReorderCalculator:
    def calculate_reorder_point(self,
        forecast: Forecast,
        lead_time: LeadTime,
        service_level: float) -> ReorderPoint:
        
        # Average demand during lead time
        avg_demand = forecast.predictions[:lead_time.days].mean()
        demand_std = forecast.predictions[:lead_time.days].std()
        
        # Safety stock for service level
        z_score = stats.norm.ppf(service_level)
        safety_stock = z_score * demand_std * np.sqrt(lead_time.days)
        
        # Reorder point
        reorder_point = (avg_demand * lead_time.days) + safety_stock
        
        # Economic order quantity
        eoq = self._calculate_eoq(avg_demand, ordering_cost, holding_cost)
        
        return ReorderPoint(
            product_id=forecast.product_id,
            reorder_point=reorder_point,
            safety_stock=safety_stock,
            order_quantity=eoq,
            next_order_date=self._predict_next_order(current_stock, reorder_point)
        )
```

TEST CASES:
1. Forecast accuracy (MAPE < 15%)
2. Seasonal pattern detection
3. Promotion impact modeling
4. Safety stock calculation accuracy
5. Reorder point trigger accuracy
6. PO generation timing
7. Multi-warehouse optimization
8. Model retraining triggers

EXPECTED DELIVERABLES:
1. Demand forecasting models (Prophet, XGBoost, LSTM)
2. Feature engineering pipeline
3. Reorder point calculator
4. Safety stock optimizer
5. Automatic PO generation service
6. Model training pipeline with MLflow
7. APIs and dashboards
8. A/B testing framework
9. Tests and documentation
```

---

## L3-P04: Fleet Management Service (Go)

```
PROJECT: OmniRoute - Fleet Management Service

CONTEXT:
Build a Fleet Management system for:
- Vehicle tracking and telematics
- Fleet maintenance scheduling
- Fuel management and monitoring
- Driver behavior analysis
- Vehicle utilization optimization
- Insurance and compliance tracking
- Vehicle allocation and assignment
- Third-party fleet integration (rentals)

TECHNICAL REQUIREMENTS:
- Language: Go 1.22+
- Telematics: OBD-II device integration via IoT gateway
- GPS: Real-time tracking with sub-minute updates
- Time-series: TimescaleDB for telemetry storage
- Events: Kafka for vehicle events
- Maps: Geofencing with PostGIS

VEHICLE TYPES:
- Motorcycle: Small package delivery, urban areas
- Three-wheeler (Keke): Medium packages, last-mile
- Van: Multiple deliveries, medium volume
- Truck: Heavy loads, B2B deliveries
- Refrigerated: Cold chain requirements

MAINTENANCE TYPES:
- Preventive: Scheduled based on mileage/time
- Predictive: Sensor-based anomaly detection
- Corrective: Breakdown repair

DOMAIN MODEL:
type Vehicle struct {
    ID                uuid.UUID
    TenantID          uuid.UUID
    RegistrationNo    string
    Type              VehicleType
    Make              string
    Model             string
    Year              int
    Capacity          VehicleCapacity
    FuelType          FuelType
    CurrentLocation   GeoPoint
    CurrentDriverID   *uuid.UUID
    Status            VehicleStatus
    OdometerReading   float64
    FuelLevel         float64
    MaintenanceStatus MaintenanceStatus
    NextServiceDue    *time.Time
    NextServiceMileage *float64
    Documents         []VehicleDocument
    Insurance         InsuranceInfo
    Sensors           []Sensor
}

TELEMATICS PROCESSING:
type TelemetryData struct {
    VehicleID    uuid.UUID
    Timestamp    time.Time
    Location     GeoPoint
    Speed        float64
    Heading      float64
    Acceleration AccelerationVector
    FuelLevel    float64
    EngineRPM    int
    OBDCodes     []string  // Diagnostic trouble codes
}

func (p *TelematicsProcessor) ProcessTelemetry(ctx context.Context,
    vehicleID uuid.UUID, data TelemetryData) error {
    // Store in TimescaleDB
    // Detect anomalies (harsh braking, speeding, idling)
    // Update vehicle state
    // Check geofence violations
    // Trigger alerts if needed
    // Publish events to Kafka
}

DRIVER BEHAVIOR SCORING:
type DriverBehavior struct {
    DriverID      uuid.UUID
    Period        DateRange
    OverallScore  float64  // 0-100
    Metrics       BehaviorMetrics
}

type BehaviorMetrics struct {
    HarshBraking     int     // Count
    HarshAcceleration int
    Speeding         float64 // % of time over limit
    Idling           float64 // % of time idling
    SeatbeltUsage    float64 // % compliance
    PhoneUsage       int     // Detected events
}

TEST CASES:
1. Real-time tracking accuracy (<30 sec delay)
2. Telematics ingestion rate (1000 vehicles/sec)
3. Predictive maintenance alerting
4. Driver behavior scoring accuracy
5. Fuel anomaly detection (theft, inefficiency)
6. Vehicle allocation optimization
7. Compliance tracking (license, insurance expiry)
8. Performance under high data volume

EXPECTED DELIVERABLES:
1. Fleet domain models
2. Telematics ingestion and processing
3. Maintenance scheduling service
4. Driver behavior scoring engine
5. Fuel management service
6. APIs and dashboards
7. Mobile driver app APIs
8. Tests and benchmarks
```

---

# Continue with remaining layers in subsequent prompts...

The document continues with:
- **Layer 4: Accessibility** (Voice, USSD, NLU, WhatsApp)
- **Layer 5: Social Commerce** (Groups, Reputation, Referrals)
- **Layer 6: Intelligence** (Market Intel, Analytics, AI Insights)
- **Layer 7: Finance** (Lending, Payments, Wallets, Credit Scoring)
- **Mobile Applications** (Customer, Worker, Admin)
- **Security & Compliance** (Auth, Privacy)
- **DevOps & MLOps** (Terraform, Helm, ML Pipelines)
- **Observability** (Full Stack)
- **Testing & Documentation**

---

# EXECUTION GUIDE

## Recommended Order

1. **Foundation (Week 1-2)**
   - D-P01: Terraform Infrastructure
   - D-P02: Helm Charts
   - S-P01: Authentication Service
   - O-P01: Observability Stack

2. **Core Commerce (Week 3-5)**
   - L1-P01: Product Catalog
   - L1-P03: Inventory Management
   - L1-P04: Customer Management
   - L1-P02: Order Management

3. **Distribution (Week 6-8)**
   - L3-P02: Warehouse Management
   - L3-P01: Route Optimization
   - L3-P04: Fleet Management
   - L3-P03: Predictive Restocking

4. **Gig Platform (Week 9-11)**
   - L2-P01: Worker Management
   - L2-P02: Task Assignment
   - L2-P03: Earnings & Benefits
   - L2-P04: Career Progression

5. **Accessibility (Week 12-14)**
   - L4-P03: Multilingual NLU
   - L4-P02: USSD Gateway
   - L4-P04: WhatsApp Integration
   - L4-P01: Voice Commerce

6. **Social & Intelligence (Week 15-17)**
   - L5-P01: Community Group Buying
   - L5-P02: Reputation Passport
   - L6-P01: Market Intelligence
   - L6-P02: Predictive Analytics

7. **Finance (Week 18-20)**
   - L7-P03: Digital Wallet
   - L7-P02: Payment Processing
   - L7-P04: Credit Scoring ML
   - L7-P01: Embedded Lending

8. **Applications (Week 21-24)**
   - M-P01: Customer Mobile App
   - M-P02: Worker Mobile App
   - M-P03: Admin Dashboard

9. **Polish (Week 25-26)**
   - L6-P03: AI Insights Generator
   - L5-P03: Referral & Affiliate
   - S-P02: Data Privacy
   - D-P03: MLOps Pipeline

---

*Document Version: 1.0 | January 2026 | BillyRonks Global Limited*
