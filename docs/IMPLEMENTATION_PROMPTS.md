# OmniRoute SCE - Comprehensive Implementation Prompts
## For AI-Assisted Development (Claude/GPT-4/Gemini/Local Models)

---

## Master Prompt: Project Context

Use this as a system prompt or prepend to all implementation prompts:

```
You are a senior software architect implementing the OmniRoute Service Creation Environment (SCE), a platform that enables non-technical users to create, extend, and customize B2B FMCG business services through visual design and AI-powered generation.

**Core Technologies:**
- Temporal (Go): Durable workflow execution
- n8n: Pre-built integrations (400+ connectors)
- AI Engine: LLMs (Claude, GPT-4) + Local models (Llama, Mistral via vLLM)
- GKE (Google Kubernetes Engine): Infrastructure

**Development Philosophy (Holy Trinity):**

1. **Extreme Programming (XP):**
   - Test-Driven Development: Write tests first, always
   - Continuous Integration: Every commit is tested
   - Simple Design: YAGNI, no premature optimization
   - Refactoring: Improve code continuously
   - Pair Programming: Two minds better than one
   - Collective Ownership: Anyone can improve any code

2. **Domain-Driven Design (DDD):**
   - Bounded Contexts: Clear service boundaries
   - Ubiquitous Language: Domain terms in code
   - Aggregates: Consistency boundaries
   - Domain Events: Eventual consistency
   - Anti-Corruption Layers: Protect domain from external systems
   - Value Objects: Immutable domain concepts

3. **Legacy Modernization:**
   - Strangler Fig Pattern: Gradual replacement
   - Event-Driven Architecture: Loose coupling
   - API Gateway: Facade for legacy systems
   - Database per Service: Avoid shared databases

**Code Quality Standards:**
- Cyclomatic Complexity: < 10 per function
- Test Coverage: > 80%
- Documentation: All public APIs documented
- Security: OWASP Top 10 compliance
- Performance: p99 latency targets met

**Output Format:**
- Always include complete, runnable code
- Include comprehensive tests (table-driven for Go)
- Add detailed comments explaining decisions
- Follow language-specific best practices
- Include error handling and edge cases
```

---

## SECTION 1: DOMAIN LAYER PROMPTS

### Prompt 1.1: Service Definition Aggregate

```
Implement the Service Definition Aggregate Root for the OmniRoute SCE.

**Domain Context:**
A ServiceDefinition represents a user-created business service that can be deployed and executed on the platform. It contains a workflow graph (nodes and edges) that defines the service's execution logic.

**Aggregate Invariants:**
1. A service cannot be published without at least one workflow node
2. A service cannot be archived if there are active workflow instances
3. Service names must be unique within a tenant
4. Version numbers must be sequential
5. Status transitions: draft → published → deprecated → archived

**Required Components:**

1. **ServiceDefinition (Aggregate Root)**
```go
type ServiceDefinition struct {
    ID            ServiceID
    TenantID      TenantID
    Name          ServiceName
    Description   string
    Category      ServiceCategory
    Workflow      WorkflowGraph
    Versions      []ServiceVersion
    ActiveVersion VersionID
    Status        ServiceStatus
    CreatedBy     UserID
    CreatedAt     time.Time
    UpdatedAt     time.Time
    PublishedAt   *time.Time
    
    // Private: domain events
    events []DomainEvent
}
```

2. **Value Objects:**
- ServiceID, TenantID, VersionID, UserID (UUID wrappers with validation)
- ServiceName (3-100 chars, alphanumeric + spaces/dashes)
- ServiceCategory (enum: automation, integration, notification, payment, fulfillment, analytics, custom)
- ServiceStatus (enum: draft, published, deprecated, archived)
- WorkflowGraph (nodes, edges, variables, triggers, error policy)
- WorkflowNode (id, type, config, position, retry policy)
- WorkflowEdge (source, target, handles, condition)

3. **Domain Events:**
- ServiceCreatedEvent
- ServicePublishedEvent
- ServiceDeprecatedEvent
- ServiceArchivedEvent
- VersionReleasedEvent

4. **Domain Methods:**
```go
func (s *ServiceDefinition) Publish() error
func (s *ServiceDefinition) Deprecate(reason string) error
func (s *ServiceDefinition) Archive() error
func (s *ServiceDefinition) AddVersion(workflow WorkflowGraph, notes string) (VersionID, error)
func (s *ServiceDefinition) UpdateWorkflow(workflow WorkflowGraph) error
func (s *ServiceDefinition) PullEvents() []DomainEvent
```

5. **Validation Rules:**
- Name: 3-100 characters, no special characters except space/dash/underscore
- Workflow: At least one node, all edges reference existing nodes
- Category: Must be from allowed list

**Test Cases (TDD - write tests first):**
1. Create service with valid data → success
2. Create service with invalid name → error
3. Publish service with workflow → success, emits event
4. Publish service without workflow → error
5. Publish archived service → error
6. Add version → increments version number, emits event
7. Deprecate published service → success, emits event
8. Archive deprecated service → success, emits event
9. Archive draft service → error (must deprecate first)

**Output Requirements:**
- domain/service_definition.go
- domain/service_definition_test.go
- domain/value_objects.go
- domain/events.go
- domain/errors.go
- All tests must pass
- 100% coverage of domain logic
```

### Prompt 1.2: Workflow Graph Value Object

```
Implement the WorkflowGraph Value Object for visual workflow representation.

**Domain Context:**
WorkflowGraph represents the visual DAG (Directed Acyclic Graph) of a service's execution flow. It is immutable and validates structural integrity.

**Structure:**
```go
type WorkflowGraph struct {
    Nodes       []WorkflowNode     
    Edges       []WorkflowEdge     
    Variables   []WorkflowVariable 
    Triggers    []WorkflowTrigger  
    ErrorPolicy ErrorHandlingPolicy
}

type WorkflowNode struct {
    ID          string
    Type        NodeType
    Label       string
    Position    Position
    Config      NodeConfig
    RetryPolicy *RetryPolicy
    Timeout     *time.Duration
}

type NodeType string
const (
    NodeTypeActivity   NodeType = "activity"      // Temporal Activity
    NodeTypeSubflow    NodeType = "subflow"       // Child Workflow
    NodeTypeAIAction   NodeType = "ai_action"     // LLM/Local Model
    NodeTypeN8N        NodeType = "n8n"           // n8n Integration
    NodeTypeDecision   NodeType = "decision"      // Conditional Branch
    NodeTypeParallel   NodeType = "parallel"      // Fork-Join
    NodeTypeWait       NodeType = "wait"          // Timer/Signal
    NodeTypeHumanTask  NodeType = "human_task"    // Human-in-loop
)

type WorkflowEdge struct {
    ID           string
    Source       string
    Target       string
    SourceHandle string
    TargetHandle string
    Condition    *EdgeCondition
}

type WorkflowVariable struct {
    Name         string
    Type         VariableType
    DefaultValue interface{}
    Required     bool
    Description  string
}

type WorkflowTrigger struct {
    Type   TriggerType
    Config TriggerConfig
}

type TriggerType string
const (
    TriggerTypeManual   TriggerType = "manual"
    TriggerTypeSchedule TriggerType = "schedule"
    TriggerTypeWebhook  TriggerType = "webhook"
    TriggerTypeEvent    TriggerType = "event"
)
```

**Validation Rules:**
1. Graph must be a valid DAG (no cycles)
2. All edge sources/targets must reference existing nodes
3. Decision nodes must have at least 2 outgoing edges
4. Parallel nodes must have matching fork/join pairs
5. At least one trigger must be defined
6. All required variables must have defaults or be inputs

**Required Methods:**
```go
func NewWorkflowGraph(nodes []WorkflowNode, edges []WorkflowEdge, ...) (*WorkflowGraph, error)
func (g WorkflowGraph) Validate() error
func (g WorkflowGraph) TopologicalSort() ([]WorkflowNode, error)
func (g WorkflowGraph) FindPath(from, to string) ([]string, error)
func (g WorkflowGraph) GetNodeByID(id string) (*WorkflowNode, bool)
func (g WorkflowGraph) GetOutgoingEdges(nodeID string) []WorkflowEdge
func (g WorkflowGraph) GetIncomingEdges(nodeID string) []WorkflowEdge
func (g WorkflowGraph) Clone() WorkflowGraph
```

**Test Cases:**
1. Valid linear workflow → success
2. Valid workflow with decision → success
3. Valid workflow with parallel fork-join → success
4. Workflow with cycle → error
5. Workflow with orphan edge → error
6. Topological sort → correct order
7. Find path between nodes → correct path
8. Clone → deep copy with no references

**Output:**
- domain/workflow_graph.go
- domain/workflow_graph_test.go
- Include graph algorithms (DFS for cycle detection, topological sort)
```

### Prompt 1.3: Repository Interfaces

```
Implement Repository interfaces following DDD patterns.

**Context:**
Repositories abstract data persistence while maintaining domain integrity. They work with aggregates and emit domain events.

**Required Interfaces:**

```go
// ServiceDefinitionRepository manages ServiceDefinition aggregates
type ServiceDefinitionRepository interface {
    // Commands
    Save(ctx context.Context, service *ServiceDefinition) error
    Delete(ctx context.Context, id ServiceID) error
    
    // Queries
    FindByID(ctx context.Context, id ServiceID) (*ServiceDefinition, error)
    FindByTenantAndName(ctx context.Context, tenantID TenantID, name ServiceName) (*ServiceDefinition, error)
    FindByTenant(ctx context.Context, tenantID TenantID, filter ServiceFilter) ([]*ServiceDefinition, error)
    CountByTenant(ctx context.Context, tenantID TenantID, filter ServiceFilter) (int64, error)
    Exists(ctx context.Context, id ServiceID) (bool, error)
    
    // Transaction support
    WithTx(ctx context.Context, fn func(ServiceDefinitionRepository) error) error
}

// ServiceFilter for querying services
type ServiceFilter struct {
    Status     []ServiceStatus
    Category   []ServiceCategory
    CreatedBy  *UserID
    Search     string
    Pagination Pagination
    SortBy     string
    SortOrder  SortOrder
}

// WorkflowInstanceRepository manages running workflow instances
type WorkflowInstanceRepository interface {
    Save(ctx context.Context, instance *WorkflowInstance) error
    FindByID(ctx context.Context, id WorkflowInstanceID) (*WorkflowInstance, error)
    FindByService(ctx context.Context, serviceID ServiceID, filter InstanceFilter) ([]*WorkflowInstance, error)
    FindActiveByService(ctx context.Context, serviceID ServiceID) ([]*WorkflowInstance, error)
    UpdateStatus(ctx context.Context, id WorkflowInstanceID, status InstanceStatus) error
}

// ActivityCatalogRepository manages available activities
type ActivityCatalogRepository interface {
    FindAll(ctx context.Context, tenantID TenantID) ([]*ActivityDefinition, error)
    FindByCategory(ctx context.Context, tenantID TenantID, category string) ([]*ActivityDefinition, error)
    FindByID(ctx context.Context, id ActivityID) (*ActivityDefinition, error)
    Register(ctx context.Context, activity *ActivityDefinition) error
    Unregister(ctx context.Context, id ActivityID) error
}
```

**Implementation Requirements:**
1. Use optimistic locking (version field)
2. Emit domain events after successful save
3. Return ErrNotFound for missing entities
4. Return ErrOptimisticLock for version conflicts
5. Support context cancellation
6. Add OpenTelemetry tracing

**Test Cases:**
1. Save new service → success
2. Save existing service → updates with version increment
3. Save with stale version → ErrOptimisticLock
4. FindByID existing → returns service
5. FindByID missing → ErrNotFound
6. FindByTenant with filter → correct results
7. Transaction rollback on error → no side effects
8. Domain events emitted after save

**Output:**
- domain/repository.go (interfaces)
- infrastructure/postgres/service_repository.go (implementation)
- infrastructure/postgres/service_repository_test.go
```

---

## SECTION 2: TEMPORAL WORKFLOW PROMPTS

### Prompt 2.1: Dynamic Workflow Executor

```
Implement the Dynamic Workflow Executor that runs user-defined services.

**Context:**
The executor interprets a WorkflowDSL (compiled from visual design) and executes it as a Temporal workflow. It must handle all node types dynamically.

**Workflow Structure:**
```go
// UserDefinedServiceWorkflow executes any user-created service
func UserDefinedServiceWorkflow(ctx workflow.Context, input *ServiceExecutionInput) (*ServiceExecutionResult, error)

type ServiceExecutionInput struct {
    ServiceID     string                 `json:"service_id"`
    VersionID     string                 `json:"version_id"`
    TenantID      string                 `json:"tenant_id"`
    WorkflowDSL   *WorkflowDSL           `json:"workflow_dsl"`
    InputData     map[string]interface{} `json:"input_data"`
    InitiatedBy   string                 `json:"initiated_by"`
    CorrelationID string                 `json:"correlation_id"`
}

type ServiceExecutionResult struct {
    Status      string                 `json:"status"`
    OutputData  map[string]interface{} `json:"output_data"`
    NodeResults map[string]*NodeResult `json:"node_results"`
    StartedAt   time.Time              `json:"started_at"`
    CompletedAt time.Time              `json:"completed_at"`
    Error       *ExecutionError        `json:"error,omitempty"`
}
```

**Node Execution Logic:**

1. **Activity Node:** Execute Temporal activity with retry policy
2. **AI Action Node:** Call AI gateway (LLM or local model)
3. **n8n Node:** Trigger n8n workflow via integration activity
4. **Decision Node:** Evaluate condition, route to appropriate branch
5. **Parallel Node:** Use workflow.Go for concurrent execution
6. **Wait Node:** Use workflow.Sleep or workflow.GetSignalChannel
7. **Human Task Node:** Signal + query pattern for human input

**Key Features:**
- Execute nodes in topological order
- Handle conditional branching (decision nodes)
- Support parallel execution with fork-join
- Implement saga pattern for compensations
- Support signals for human-in-the-loop
- Expose queries for workflow state

**Implementation:**
```go
func UserDefinedServiceWorkflow(ctx workflow.Context, input *ServiceExecutionInput) (*ServiceExecutionResult, error) {
    logger := workflow.GetLogger(ctx)
    logger.Info("Starting service execution", "serviceID", input.ServiceID)
    
    // Initialize execution context
    execCtx := &ExecutionContext{
        Variables:   make(map[string]interface{}),
        NodeResults: make(map[string]*NodeResult),
    }
    
    // Copy input data to variables
    for k, v := range input.InputData {
        execCtx.Variables[k] = v
    }
    
    // Get execution order (topological sort)
    executionOrder, err := getExecutionOrder(input.WorkflowDSL)
    if err != nil {
        return nil, err
    }
    
    // Execute nodes
    for _, nodeID := range executionOrder {
        node := input.WorkflowDSL.GetNode(nodeID)
        
        result, err := executeNode(ctx, node, execCtx, input)
        if err != nil {
            if input.WorkflowDSL.ErrorPolicy.OnError == "compensate" {
                return compensate(ctx, execCtx, err)
            }
            return nil, err
        }
        
        execCtx.NodeResults[nodeID] = result
        
        // Update variables with output
        for k, v := range result.Output {
            execCtx.Variables[k] = v
        }
    }
    
    return &ServiceExecutionResult{
        Status:      "completed",
        OutputData:  execCtx.Variables,
        NodeResults: execCtx.NodeResults,
        CompletedAt: workflow.Now(ctx),
    }, nil
}

func executeNode(ctx workflow.Context, node *WorkflowNode, execCtx *ExecutionContext, input *ServiceExecutionInput) (*NodeResult, error) {
    switch node.Type {
    case NodeTypeActivity:
        return executeActivityNode(ctx, node, execCtx)
    case NodeTypeAIAction:
        return executeAIActionNode(ctx, node, execCtx)
    case NodeTypeN8N:
        return executeN8NNode(ctx, node, execCtx)
    case NodeTypeDecision:
        return executeDecisionNode(ctx, node, execCtx, input.WorkflowDSL)
    case NodeTypeParallel:
        return executeParallelNode(ctx, node, execCtx, input)
    case NodeTypeWait:
        return executeWaitNode(ctx, node, execCtx)
    case NodeTypeHumanTask:
        return executeHumanTaskNode(ctx, node, execCtx)
    default:
        return nil, fmt.Errorf("unknown node type: %s", node.Type)
    }
}
```

**Test Cases:**
1. Linear workflow (3 activities) → all execute in order
2. Decision workflow (if-else) → correct branch taken
3. Parallel workflow (fork-join) → concurrent execution
4. Workflow with signal wait → pauses until signal
5. Workflow with human task → waits for approval
6. Workflow with error → compensation executed
7. Query workflow state → returns current status

**Output:**
- workflows/user_service_workflow.go
- workflows/node_executors.go
- workflows/saga.go
- workflows/user_service_workflow_test.go
```

### Prompt 2.2: Core Activities

```
Implement Core Activities for the Temporal workers.

**Context:**
Activities are the building blocks of workflows. They perform actual work like API calls, database operations, and notifications.

**Activity Categories:**

1. **HTTP Activities:**
```go
type HTTPActivityInput struct {
    Method      string            `json:"method"`
    URL         string            `json:"url"`
    Headers     map[string]string `json:"headers"`
    Body        interface{}       `json:"body"`
    Timeout     time.Duration     `json:"timeout"`
    RetryPolicy *RetryPolicy      `json:"retry_policy"`
}

func (a *CoreActivities) HTTPCall(ctx context.Context, input *HTTPActivityInput) (*HTTPActivityResult, error)
```

2. **Database Activities:**
```go
type QueryActivityInput struct {
    Connection string                 `json:"connection"` // Reference to connection config
    Query      string                 `json:"query"`
    Params     []interface{}          `json:"params"`
    Timeout    time.Duration          `json:"timeout"`
}

func (a *CoreActivities) ExecuteQuery(ctx context.Context, input *QueryActivityInput) (*QueryActivityResult, error)
func (a *CoreActivities) ExecuteCommand(ctx context.Context, input *CommandActivityInput) (*CommandActivityResult, error)
```

3. **Notification Activities:**
```go
type SendNotificationInput struct {
    TenantID  string            `json:"tenant_id"`
    Channel   string            `json:"channel"` // sms, email, push, whatsapp
    Recipient string            `json:"recipient"`
    Template  string            `json:"template"`
    Variables map[string]string `json:"variables"`
}

func (a *CoreActivities) SendNotification(ctx context.Context, input *SendNotificationInput) (*SendNotificationResult, error)
```

4. **Transform Activities:**
```go
func (a *CoreActivities) TransformJSON(ctx context.Context, input *TransformInput) (*TransformResult, error)
func (a *CoreActivities) ValidateSchema(ctx context.Context, input *ValidateInput) (*ValidateResult, error)
func (a *CoreActivities) MapFields(ctx context.Context, input *MapFieldsInput) (*MapFieldsResult, error)
```

5. **Storage Activities:**
```go
func (a *CoreActivities) UploadFile(ctx context.Context, input *UploadInput) (*UploadResult, error)
func (a *CoreActivities) DownloadFile(ctx context.Context, input *DownloadInput) (*DownloadResult, error)
func (a *CoreActivities) GenerateSignedURL(ctx context.Context, input *SignedURLInput) (*SignedURLResult, error)
```

**Implementation Requirements:**
- All activities must be idempotent
- Use heartbeats for long-running operations (>10s)
- Include OpenTelemetry spans
- Implement circuit breakers for external calls
- Log all activity executions
- Return structured errors

**Activity Registration:**
```go
type CoreActivities struct {
    httpClient  *http.Client
    dbPool      *pgxpool.Pool
    notifClient notification.Client
    storage     storage.Client
    breakers    map[string]*gobreaker.CircuitBreaker
}

func NewCoreActivities(deps Dependencies) *CoreActivities {
    return &CoreActivities{
        httpClient:  deps.HTTPClient,
        dbPool:      deps.DBPool,
        notifClient: deps.NotificationClient,
        storage:     deps.StorageClient,
        breakers:    initCircuitBreakers(),
    }
}
```

**Test Cases:**
1. HTTP call success → returns response
2. HTTP call timeout → returns error with context
3. HTTP call with retry → retries on 5xx
4. Database query → returns rows
5. Send notification → success
6. Upload file with heartbeat → completes
7. Circuit breaker opens after failures

**Output:**
- activities/core_activities.go
- activities/http.go
- activities/database.go
- activities/notification.go
- activities/storage.go
- activities/core_activities_test.go
```

### Prompt 2.3: Integration Activities (n8n + AI)

```
Implement Integration Activities for n8n and AI operations.

**n8n Integration:**
```go
type IntegrationActivities struct {
    n8nGateway  *n8n.Gateway
    aiGateway   *ai.Gateway
    breaker     *gobreaker.CircuitBreaker
}

// ExecuteN8NWorkflow triggers an n8n workflow
type ExecuteN8NWorkflowInput struct {
    TenantID          string                 `json:"tenant_id"`
    WorkflowID        string                 `json:"workflow_id"`
    WebhookPath       string                 `json:"webhook_path,omitempty"`
    InputData         map[string]interface{} `json:"input_data"`
    WaitForCompletion bool                   `json:"wait_for_completion"`
    Timeout           time.Duration          `json:"timeout"`
}

func (a *IntegrationActivities) ExecuteN8NWorkflow(ctx context.Context, input *ExecuteN8NWorkflowInput) (*N8NWorkflowResult, error)

// ListN8NWorkflows returns available n8n workflows for tenant
func (a *IntegrationActivities) ListN8NWorkflows(ctx context.Context, tenantID string) ([]N8NWorkflowMeta, error)
```

**AI Integration:**
```go
// CallLLM invokes a cloud LLM (Claude, GPT-4, Gemini)
type CallLLMInput struct {
    Provider    string            `json:"provider"` // anthropic, openai, google
    Model       string            `json:"model"`
    Messages    []Message         `json:"messages"`
    System      string            `json:"system,omitempty"`
    Temperature float64           `json:"temperature"`
    MaxTokens   int               `json:"max_tokens"`
    Tools       []Tool            `json:"tools,omitempty"`
}

func (a *IntegrationActivities) CallLLM(ctx context.Context, input *CallLLMInput) (*LLMResult, error)

// LocalModelInference invokes local model via vLLM
type LocalModelInferenceInput struct {
    Model       string    `json:"model"` // llama-3.3-70b, mistral-large, qwen-coder
    Messages    []Message `json:"messages"`
    Temperature float64   `json:"temperature"`
    MaxTokens   int       `json:"max_tokens"`
}

func (a *IntegrationActivities) LocalModelInference(ctx context.Context, input *LocalModelInferenceInput) (*LLMResult, error)

// GenerateEmbedding creates vector embeddings
type GenerateEmbeddingInput struct {
    Model string   `json:"model"`
    Texts []string `json:"texts"`
}

func (a *IntegrationActivities) GenerateEmbedding(ctx context.Context, input *GenerateEmbeddingInput) (*EmbeddingResult, error)

// StructuredOutput generates structured data from LLM
type StructuredOutputInput struct {
    Provider string      `json:"provider"`
    Model    string      `json:"model"`
    Prompt   string      `json:"prompt"`
    Schema   interface{} `json:"schema"` // JSON Schema
}

func (a *IntegrationActivities) StructuredOutput(ctx context.Context, input *StructuredOutputInput) (*StructuredOutputResult, error)
```

**Implementation Details:**

1. **n8n Integration:**
   - Support webhook trigger and API execution
   - Wait for completion with polling
   - Handle n8n execution states
   - Map n8n data to workflow variables

2. **AI Integration:**
   - Route to appropriate provider based on input
   - Implement fallback (cloud → local)
   - Track token usage for billing
   - Cache common prompts
   - Handle streaming responses

3. **Error Handling:**
   - Retry on transient errors (rate limits, timeouts)
   - Circuit breaker for repeated failures
   - Structured error messages

**Test Cases:**
1. Execute n8n workflow via webhook → success
2. Execute n8n workflow via API → success
3. Wait for n8n completion → returns result
4. Call Claude → returns completion
5. Call GPT-4 → returns completion
6. Local inference → returns completion
7. Fallback from cloud to local → success
8. Generate embeddings → returns vectors

**Output:**
- activities/integration_activities.go
- activities/n8n.go
- activities/ai.go
- activities/integration_activities_test.go
```

---

## SECTION 3: AI SERVICE GENERATION PROMPTS

### Prompt 3.1: Prompt-to-Service Generator

```
Implement the Prompt-to-Service generation engine.

**Context:**
Users describe services in natural language, and the AI generates complete service definitions with workflows.

**API Endpoint:**
```python
@app.post("/api/v1/generate/service")
async def generate_service(request: ServiceGenerationRequest) -> ServiceGenerationResponse:
    """
    Generate a complete service definition from natural language.
    
    Example prompts:
    - "Create a service that monitors orders and sends SMS when status changes"
    - "Build a workflow that fetches data from our ERP, processes it with AI, and updates the CRM"
    - "Design a payment collection service with mobile money, bank transfer fallback"
    """
```

**Generation Pipeline:**

1. **Prompt Analysis:**
   - Extract intents (create, modify, integrate)
   - Identify entities (services, data sources, actions)
   - Determine workflow pattern (linear, branching, parallel)

2. **Context Enrichment:**
   - Load available activities
   - Load available n8n workflows
   - Load tenant-specific configurations
   - Load similar service templates

3. **Service Generation:**
   - Generate service metadata (name, description, category)
   - Generate workflow nodes
   - Generate workflow edges
   - Generate triggers
   - Generate error handling

4. **Validation:**
   - Validate workflow structure (DAG, no orphans)
   - Validate node configurations
   - Validate referenced activities exist
   - Validate permissions

**System Prompt Template:**
```python
SYSTEM_PROMPT = """
You are a service architect for OmniRoute, a B2B FMCG platform.
You help users create business services that orchestrate:
- Temporal workflows (durable, fault-tolerant)
- n8n integrations (400+ connectors)
- AI capabilities (LLM and local models)

Available node types:
1. activity: Execute a Temporal activity (HTTP, database, notification)
2. ai_action: Call an LLM or local model
3. n8n: Execute an n8n workflow
4. decision: Conditional branching (if/else)
5. parallel: Concurrent execution (fork/join)
6. wait: Timer or signal wait
7. human_task: Human approval/input

Available activities:
{activities_list}

Available n8n workflows:
{n8n_workflows_list}

When generating a service, output a JSON object with this structure:
{
  "name": "ServiceName",
  "description": "What the service does",
  "category": "automation|integration|notification|payment|fulfillment|analytics",
  "workflow": {
    "nodes": [
      {
        "id": "unique_id",
        "type": "activity|ai_action|n8n|decision|parallel|wait|human_task",
        "label": "Human readable label",
        "position": {"x": 100, "y": 100},
        "config": {
          // Type-specific configuration
        },
        "retry_policy": {
          "max_attempts": 3,
          "initial_interval": "1s",
          "backoff_coefficient": 2.0
        }
      }
    ],
    "edges": [
      {
        "id": "edge_id",
        "source": "node_id",
        "target": "node_id",
        "condition": null  // Optional condition for decision nodes
      }
    ],
    "triggers": [
      {
        "type": "manual|schedule|webhook|event",
        "config": {}
      }
    ],
    "variables": [
      {
        "name": "variable_name",
        "type": "string|number|boolean|object|array",
        "required": true,
        "description": "What this variable is for"
      }
    ],
    "error_policy": {
      "on_error": "fail|compensate|ignore",
      "compensation_workflow": null
    }
  },
  "inputs": [
    {
      "name": "input_name",
      "type": "string",
      "required": true,
      "description": "Description"
    }
  ],
  "outputs": [
    {
      "name": "output_name",
      "type": "string",
      "description": "Description"
    }
  ]
}

After the JSON, provide:
1. A brief explanation of what the service does
2. Any assumptions made
3. Suggestions for enhancements

IMPORTANT: 
- Always use descriptive node IDs
- Position nodes for visual clarity (left-to-right flow)
- Include appropriate retry policies
- Add meaningful labels
"""
```

**Implementation:**
```python
class ServiceGenerator:
    def __init__(self, ai_gateway: AIGateway, activity_catalog: ActivityCatalog):
        self.ai_gateway = ai_gateway
        self.activity_catalog = activity_catalog
    
    async def generate(
        self,
        tenant_id: str,
        prompt: str,
        context: Optional[Dict[str, Any]] = None,
        use_local_model: bool = False,
    ) -> ServiceGenerationResult:
        # 1. Load context
        activities = await self.activity_catalog.get_for_tenant(tenant_id)
        n8n_workflows = await self.n8n_client.list_workflows(tenant_id)
        
        # 2. Build system prompt
        system_prompt = self._build_system_prompt(activities, n8n_workflows)
        
        # 3. Build user prompt
        user_prompt = self._build_user_prompt(prompt, context)
        
        # 4. Generate
        if use_local_model:
            response = await self.ai_gateway.local_inference(
                model="llama-3.3-70b-instruct",
                system=system_prompt,
                messages=[{"role": "user", "content": user_prompt}],
                temperature=0.7,
                max_tokens=4000,
            )
        else:
            response = await self.ai_gateway.call_llm(
                provider="anthropic",
                model="claude-sonnet-4-20250514",
                system=system_prompt,
                messages=[{"role": "user", "content": user_prompt}],
                temperature=0.7,
                max_tokens=4000,
            )
        
        # 5. Parse response
        parsed = self._parse_response(response.content)
        
        # 6. Validate
        errors = await self._validate(parsed.service_definition, tenant_id)
        
        if errors:
            return ServiceGenerationResult(
                status="error",
                errors=errors,
                suggestions=parsed.suggestions,
            )
        
        return ServiceGenerationResult(
            status="success",
            service_definition=parsed.service_definition,
            explanation=parsed.explanation,
            suggestions=parsed.suggestions,
            tokens_used=response.usage.total_tokens,
            model_used=response.model,
        )
    
    def _parse_response(self, content: str) -> ParsedResponse:
        # Extract JSON from response
        json_match = re.search(r'\{[\s\S]*\}', content)
        if not json_match:
            raise ParseError("No JSON found in response")
        
        service_def = json.loads(json_match.group())
        
        # Extract explanation (text after JSON)
        explanation = content[json_match.end():].strip()
        
        # Extract suggestions
        suggestions = self._extract_suggestions(explanation)
        
        return ParsedResponse(
            service_definition=service_def,
            explanation=explanation,
            suggestions=suggestions,
        )
```

**Test Cases:**
1. Simple linear workflow → generates correct nodes and edges
2. Workflow with decision → generates conditional branches
3. Workflow with AI action → includes AI node configuration
4. Workflow with n8n → references valid n8n workflow
5. Invalid prompt → returns helpful error
6. Local model generation → works correctly
7. Validation catches invalid references

**Output:**
- app/services/generator.py
- app/prompts/service_generation.py
- app/parsers/response_parser.py
- tests/test_generator.py
```

### Prompt 3.2: Workflow Enhancement

```
Implement the Workflow Enhancement feature.

**Context:**
Users can enhance existing workflows using natural language instructions.

**API Endpoint:**
```python
@app.post("/api/v1/enhance/workflow")
async def enhance_workflow(
    tenant_id: str,
    workflow_dsl: Dict[str, Any],
    enhancement_prompt: str,
) -> WorkflowEnhancementResult:
    """
    Enhance an existing workflow based on natural language.
    
    Example prompts:
    - "Add retry logic with exponential backoff to all HTTP calls"
    - "Insert a notification step after the payment"
    - "Add error handling that creates a support ticket"
    - "Make the AI step use a local model instead"
    """
```

**Enhancement Types:**

1. **Node Modifications:**
   - Add retry policy
   - Change timeout
   - Update configuration
   - Change node type

2. **Structural Changes:**
   - Insert new node
   - Remove node
   - Add parallel branch
   - Add decision branch

3. **Error Handling:**
   - Add try-catch pattern
   - Add compensation
   - Add notification on failure

4. **AI-Specific:**
   - Switch provider
   - Update prompt
   - Add structured output

**System Prompt:**
```python
ENHANCEMENT_PROMPT = """
You are a workflow enhancement expert.

Given an existing workflow DSL and an enhancement request, modify the workflow.

Rules:
1. Preserve existing functionality unless explicitly asked to change it
2. Maintain node ID consistency (don't rename existing IDs)
3. Update edges when adding/removing nodes
4. Keep positions logical (new nodes flow left-to-right)
5. Add appropriate retry policies for new nodes

Current workflow:
{workflow_json}

Enhancement request: {enhancement_prompt}

Output the COMPLETE modified workflow DSL (not just the changes).
After the JSON, explain what changes were made.
"""
```

**Implementation:**
```python
class WorkflowEnhancer:
    async def enhance(
        self,
        tenant_id: str,
        workflow_dsl: Dict[str, Any],
        enhancement_prompt: str,
    ) -> WorkflowEnhancementResult:
        # 1. Validate current workflow
        validation_errors = self._validate_workflow(workflow_dsl)
        if validation_errors:
            raise InvalidWorkflowError(validation_errors)
        
        # 2. Build enhancement prompt
        system_prompt = self._build_system_prompt()
        user_prompt = ENHANCEMENT_PROMPT.format(
            workflow_json=json.dumps(workflow_dsl, indent=2),
            enhancement_prompt=enhancement_prompt,
        )
        
        # 3. Generate enhancement
        response = await self.ai_gateway.call_llm(
            provider="anthropic",
            model="claude-sonnet-4-20250514",
            system=system_prompt,
            messages=[{"role": "user", "content": user_prompt}],
            temperature=0.3,  # Lower temperature for more predictable changes
            max_tokens=8000,
        )
        
        # 4. Parse enhanced workflow
        enhanced = self._parse_response(response.content)
        
        # 5. Validate enhanced workflow
        validation_errors = self._validate_workflow(enhanced.workflow)
        if validation_errors:
            return WorkflowEnhancementResult(
                status="error",
                errors=validation_errors,
            )
        
        # 6. Calculate diff
        diff = self._calculate_diff(workflow_dsl, enhanced.workflow)
        
        return WorkflowEnhancementResult(
            status="success",
            enhanced_workflow=enhanced.workflow,
            changes=diff,
            explanation=enhanced.explanation,
        )
    
    def _calculate_diff(
        self,
        original: Dict[str, Any],
        enhanced: Dict[str, Any],
    ) -> WorkflowDiff:
        # Compare nodes
        original_nodes = {n["id"]: n for n in original.get("nodes", [])}
        enhanced_nodes = {n["id"]: n for n in enhanced.get("nodes", [])}
        
        added_nodes = [n for id, n in enhanced_nodes.items() if id not in original_nodes]
        removed_nodes = [n for id, n in original_nodes.items() if id not in enhanced_nodes]
        modified_nodes = [
            {"before": original_nodes[id], "after": enhanced_nodes[id]}
            for id in original_nodes.keys() & enhanced_nodes.keys()
            if original_nodes[id] != enhanced_nodes[id]
        ]
        
        # Compare edges
        original_edges = {(e["source"], e["target"]): e for e in original.get("edges", [])}
        enhanced_edges = {(e["source"], e["target"]): e for e in enhanced.get("edges", [])}
        
        added_edges = [e for k, e in enhanced_edges.items() if k not in original_edges]
        removed_edges = [e for k, e in original_edges.items() if k not in enhanced_edges]
        
        return WorkflowDiff(
            added_nodes=added_nodes,
            removed_nodes=removed_nodes,
            modified_nodes=modified_nodes,
            added_edges=added_edges,
            removed_edges=removed_edges,
        )
```

**Test Cases:**
1. Add retry policy → all applicable nodes updated
2. Insert notification → node added, edges updated
3. Add error handling → compensation nodes added
4. Switch AI provider → config updated
5. Invalid enhancement → returns helpful error
6. Diff calculation → shows correct changes

**Output:**
- app/services/enhancer.py
- app/parsers/diff_calculator.py
- tests/test_enhancer.py
```

---

## SECTION 4: FRONTEND PROMPTS

### Prompt 4.1: Service Designer Canvas

```
Implement the Service Designer canvas using React Flow.

**Context:**
Non-technical users design services by dragging nodes onto a canvas and connecting them visually.

**Component Structure:**
```
src/
├── components/
│   └── ServiceDesigner/
│       ├── ServiceDesigner.tsx          # Main component
│       ├── Canvas.tsx                   # React Flow wrapper
│       ├── nodes/                       # Custom node components
│       │   ├── BaseNode.tsx
│       │   ├── ActivityNode.tsx
│       │   ├── AIActionNode.tsx
│       │   ├── N8NNode.tsx
│       │   ├── DecisionNode.tsx
│       │   ├── ParallelNode.tsx
│       │   ├── WaitNode.tsx
│       │   └── HumanTaskNode.tsx
│       ├── NodePalette.tsx              # Draggable node palette
│       ├── PropertyPanel.tsx            # Node configuration panel
│       ├── AIAssistant.tsx              # AI chat interface
│       └── ValidationPanel.tsx          # Validation errors
├── stores/
│   └── serviceStore.ts                  # Zustand store
├── lib/
│   └── workflowCompiler.ts              # DSL compiler
└── types/
    └── workflow.ts                       # TypeScript types
```

**Node Component Template:**
```tsx
// components/ServiceDesigner/nodes/BaseNode.tsx
import { memo } from 'react';
import { Handle, Position, NodeProps } from 'reactflow';

interface BaseNodeData {
  label: string;
  config: Record<string, unknown>;
  isValid: boolean;
  validationErrors?: string[];
}

export const BaseNode = memo(({ data, selected }: NodeProps<BaseNodeData>) => {
  return (
    <div
      className={cn(
        'px-4 py-2 rounded-lg border-2 shadow-sm min-w-[150px]',
        selected ? 'border-blue-500' : 'border-gray-200',
        !data.isValid && 'border-red-500'
      )}
    >
      <Handle type="target" position={Position.Top} />
      
      <div className="flex items-center gap-2">
        <span className="text-lg">{getNodeIcon(data)}</span>
        <span className="font-medium text-sm">{data.label}</span>
      </div>
      
      {data.validationErrors && data.validationErrors.length > 0 && (
        <div className="text-xs text-red-500 mt-1">
          {data.validationErrors[0]}
        </div>
      )}
      
      <Handle type="source" position={Position.Bottom} />
    </div>
  );
});
```

**Store Implementation:**
```typescript
// stores/serviceStore.ts
import { create } from 'zustand';
import { Node, Edge, applyNodeChanges, applyEdgeChanges } from 'reactflow';

interface ServiceState {
  // Service metadata
  serviceId: string | null;
  serviceName: string;
  serviceDescription: string;
  serviceCategory: string;
  
  // Workflow
  nodes: Node[];
  edges: Edge[];
  
  // UI state
  selectedNodeId: string | null;
  isAIAssistantOpen: boolean;
  validationErrors: ValidationError[];
  
  // Actions
  setServiceMeta: (meta: Partial<ServiceMeta>) => void;
  addNode: (node: Node) => void;
  updateNode: (nodeId: string, data: Partial<NodeData>) => void;
  removeNode: (nodeId: string) => void;
  addEdge: (edge: Edge) => void;
  removeEdge: (edgeId: string) => void;
  selectNode: (nodeId: string | null) => void;
  setAIAssistantOpen: (open: boolean) => void;
  validate: () => Promise<ValidationError[]>;
  compile: () => WorkflowDSL;
  loadFromDSL: (dsl: WorkflowDSL) => void;
  reset: () => void;
}

export const useServiceStore = create<ServiceState>((set, get) => ({
  // Initial state
  serviceId: null,
  serviceName: '',
  serviceDescription: '',
  serviceCategory: 'automation',
  nodes: [],
  edges: [],
  selectedNodeId: null,
  isAIAssistantOpen: false,
  validationErrors: [],
  
  // Actions
  addNode: (node) => {
    set((state) => ({
      nodes: [...state.nodes, node],
    }));
  },
  
  updateNode: (nodeId, data) => {
    set((state) => ({
      nodes: state.nodes.map((node) =>
        node.id === nodeId
          ? { ...node, data: { ...node.data, ...data } }
          : node
      ),
    }));
  },
  
  compile: () => {
    const { nodes, edges, serviceName, serviceDescription, serviceCategory } = get();
    return compileToWorkflowDSL(nodes, edges, {
      name: serviceName,
      description: serviceDescription,
      category: serviceCategory,
    });
  },
  
  validate: async () => {
    const dsl = get().compile();
    const errors = await validateWorkflow(dsl);
    set({ validationErrors: errors });
    return errors;
  },
}));
```

**Workflow Compiler:**
```typescript
// lib/workflowCompiler.ts
export function compileToWorkflowDSL(
  nodes: Node[],
  edges: Edge[],
  meta: ServiceMeta
): WorkflowDSL {
  return {
    name: meta.name,
    description: meta.description,
    category: meta.category,
    workflow: {
      nodes: nodes.map((node) => ({
        id: node.id,
        type: node.type as NodeType,
        label: node.data.label,
        position: node.position,
        config: node.data.config,
        retry_policy: node.data.retryPolicy,
      })),
      edges: edges.map((edge) => ({
        id: edge.id,
        source: edge.source,
        target: edge.target,
        source_handle: edge.sourceHandle,
        target_handle: edge.targetHandle,
        condition: edge.data?.condition,
      })),
      triggers: [{ type: 'manual', config: {} }],
      variables: extractVariables(nodes),
      error_policy: { on_error: 'fail' },
    },
  };
}

export async function validateWorkflow(dsl: WorkflowDSL): Promise<ValidationError[]> {
  const errors: ValidationError[] = [];
  
  // Check for cycles
  if (hasCycle(dsl.workflow.nodes, dsl.workflow.edges)) {
    errors.push({
      type: 'structure',
      message: 'Workflow contains a cycle',
    });
  }
  
  // Check for orphan edges
  const nodeIds = new Set(dsl.workflow.nodes.map((n) => n.id));
  for (const edge of dsl.workflow.edges) {
    if (!nodeIds.has(edge.source)) {
      errors.push({
        type: 'edge',
        edgeId: edge.id,
        message: `Edge source "${edge.source}" does not exist`,
      });
    }
    if (!nodeIds.has(edge.target)) {
      errors.push({
        type: 'edge',
        edgeId: edge.id,
        message: `Edge target "${edge.target}" does not exist`,
      });
    }
  }
  
  // Check decision nodes have multiple outputs
  for (const node of dsl.workflow.nodes) {
    if (node.type === 'decision') {
      const outgoingEdges = dsl.workflow.edges.filter((e) => e.source === node.id);
      if (outgoingEdges.length < 2) {
        errors.push({
          type: 'node',
          nodeId: node.id,
          message: 'Decision node must have at least 2 outgoing edges',
        });
      }
    }
  }
  
  return errors;
}
```

**Test Cases:**
1. Add node via drag-and-drop → node appears on canvas
2. Connect nodes → edge created
3. Delete node → node and connected edges removed
4. Update node config → property panel reflects changes
5. Validate valid workflow → no errors
6. Validate invalid workflow → shows errors
7. Compile to DSL → correct structure
8. Load from DSL → nodes and edges rendered

**Output:**
- All component files listed above
- Comprehensive tests using Vitest and React Testing Library
```

---

## SECTION 5: INFRASTRUCTURE PROMPTS

### Prompt 5.1: GKE Terraform Configuration

```
Implement GKE infrastructure using Terraform.

**Requirements:**
1. GKE Autopilot cluster
2. Cloud SQL PostgreSQL (HA)
3. Memorystore Redis (HA)
4. Cloud Storage for model cache
5. Secret Manager for credentials
6. Cloud Armor for security
7. Cloud CDN for static assets

**File Structure:**
```
terraform/
├── environments/
│   ├── dev/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── terraform.tfvars
│   ├── staging/
│   └── production/
├── modules/
│   ├── gke/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── outputs.tf
│   ├── cloudsql/
│   ├── redis/
│   ├── networking/
│   ├── security/
│   └── observability/
└── scripts/
    └── deploy.sh
```

**GKE Module:**
```hcl
# modules/gke/main.tf
resource "google_container_cluster" "primary" {
  name     = var.cluster_name
  location = var.region
  
  # Use Autopilot for serverless K8s
  enable_autopilot = true
  
  # Network configuration
  network    = var.network_id
  subnetwork = var.subnet_id
  
  ip_allocation_policy {
    cluster_secondary_range_name  = var.pods_range_name
    services_secondary_range_name = var.services_range_name
  }
  
  # Workload Identity
  workload_identity_config {
    workload_pool = "${var.project_id}.svc.id.goog"
  }
  
  # Binary Authorization
  binary_authorization {
    evaluation_mode = "PROJECT_SINGLETON_POLICY_ENFORCE"
  }
  
  # Logging and Monitoring
  logging_config {
    enable_components = ["SYSTEM_COMPONENTS", "WORKLOADS"]
  }
  
  monitoring_config {
    enable_components = ["SYSTEM_COMPONENTS", "WORKLOADS"]
    managed_prometheus {
      enabled = true
    }
  }
  
  # Maintenance window (3 AM UTC)
  maintenance_policy {
    daily_maintenance_window {
      start_time = "03:00"
    }
  }
  
  # Release channel
  release_channel {
    channel = "REGULAR"
  }
  
  # Addons
  addons_config {
    http_load_balancing {
      disabled = false
    }
    horizontal_pod_autoscaling {
      disabled = false
    }
    gce_persistent_disk_csi_driver_config {
      enabled = true
    }
  }
  
  # Private cluster
  private_cluster_config {
    enable_private_nodes    = true
    enable_private_endpoint = false
    master_ipv4_cidr_block  = var.master_ipv4_cidr_block
  }
  
  # Master authorized networks
  master_authorized_networks_config {
    cidr_blocks {
      cidr_block   = var.authorized_network
      display_name = "Authorized Network"
    }
  }
}

# Node pool for GPU workloads (if needed beyond Autopilot)
resource "google_container_node_pool" "gpu" {
  count = var.enable_gpu_pool ? 1 : 0
  
  name       = "gpu-pool"
  location   = var.region
  cluster    = google_container_cluster.primary.name
  
  initial_node_count = 0
  
  autoscaling {
    min_node_count = 0
    max_node_count = var.gpu_max_nodes
  }
  
  node_config {
    machine_type = "a2-highgpu-4g"
    
    guest_accelerator {
      type  = "nvidia-tesla-a100"
      count = 4
      gpu_driver_installation_config {
        gpu_driver_version = "LATEST"
      }
    }
    
    disk_size_gb = 500
    disk_type    = "pd-ssd"
    
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
    
    labels = {
      workload = "ai-inference"
    }
    
    taint {
      key    = "nvidia.com/gpu"
      value  = "present"
      effect = "NO_SCHEDULE"
    }
    
    # Workload Identity
    workload_metadata_config {
      mode = "GKE_METADATA"
    }
  }
  
  management {
    auto_repair  = true
    auto_upgrade = true
  }
}
```

**Cloud SQL Module:**
```hcl
# modules/cloudsql/main.tf
resource "google_sql_database_instance" "primary" {
  name             = var.instance_name
  database_version = "POSTGRES_15"
  region           = var.region
  
  deletion_protection = var.deletion_protection
  
  settings {
    tier              = var.machine_type
    availability_type = "REGIONAL"
    disk_autoresize   = true
    disk_size         = var.disk_size
    disk_type         = "PD_SSD"
    
    backup_configuration {
      enabled                        = true
      point_in_time_recovery_enabled = true
      start_time                     = "03:00"
      
      backup_retention_settings {
        retained_backups = 30
        retention_unit   = "COUNT"
      }
    }
    
    ip_configuration {
      ipv4_enabled    = false
      private_network = var.network_id
      require_ssl     = true
    }
    
    maintenance_window {
      day          = 7  # Sunday
      hour         = 3
      update_track = "stable"
    }
    
    database_flags {
      name  = "max_connections"
      value = "500"
    }
    
    database_flags {
      name  = "work_mem"
      value = "64MB"
    }
    
    database_flags {
      name  = "shared_preload_libraries"
      value = "pg_stat_statements"
    }
    
    insights_config {
      query_insights_enabled  = true
      query_string_length     = 2048
      record_application_tags = true
      record_client_address   = true
    }
  }
}

resource "google_sql_database" "databases" {
  for_each = toset(var.databases)
  
  name     = each.value
  instance = google_sql_database_instance.primary.name
}

resource "google_sql_user" "users" {
  for_each = var.users
  
  name     = each.key
  instance = google_sql_database_instance.primary.name
  password = each.value.password
}
```

**Output:**
- All Terraform modules
- Environment configurations
- Deploy script with proper ordering
- Documentation
```

### Prompt 5.2: Helm Charts

```
Implement Helm charts for SCE deployment.

**Chart Structure:**
```
helm/
└── omniroute-sce/
    ├── Chart.yaml
    ├── values.yaml
    ├── values-dev.yaml
    ├── values-prod.yaml
    ├── templates/
    │   ├── _helpers.tpl
    │   ├── NOTES.txt
    │   ├── configmap.yaml
    │   ├── secret.yaml
    │   ├── service-registry/
    │   │   ├── deployment.yaml
    │   │   ├── service.yaml
    │   │   └── hpa.yaml
    │   ├── workflow-compiler/
    │   ├── ai-gateway/
    │   ├── temporal/
    │   ├── n8n/
    │   ├── frontend/
    │   └── ingress.yaml
    └── charts/
        └── temporal/
```

**Main values.yaml:**
```yaml
global:
  namespace: omniroute-sce
  
  image:
    registry: gcr.io/omniroute-platform
    pullPolicy: IfNotPresent
    pullSecrets: []
  
  environment: production
  
  # OpenTelemetry
  otel:
    enabled: true
    endpoint: "http://otel-collector.observability:4317"
    samplingRatio: 0.1
  
  # PostgreSQL
  postgres:
    host: ""  # Set in environment values
    port: 5432
    sslMode: require
    existingSecret: postgres-credentials
  
  # Redis
  redis:
    host: ""  # Set in environment values
    port: 6379
    tls: true
    existingSecret: redis-credentials

# Service Registry
serviceRegistry:
  enabled: true
  replicaCount: 3
  
  image:
    repository: sce-service-registry
    tag: latest
  
  resources:
    requests:
      memory: "256Mi"
      cpu: "100m"
    limits:
      memory: "1Gi"
      cpu: "500m"
  
  autoscaling:
    enabled: true
    minReplicas: 2
    maxReplicas: 10
    targetCPUUtilization: 70
    targetMemoryUtilization: 80
  
  service:
    type: ClusterIP
    port: 8080
  
  healthcheck:
    path: /health
    port: 8080
  
  env: []
  # Additional environment variables

# Workflow Compiler (Rust)
workflowCompiler:
  enabled: true
  replicaCount: 3
  
  image:
    repository: sce-workflow-compiler
    tag: latest
  
  resources:
    requests:
      memory: "256Mi"
      cpu: "100m"
    limits:
      memory: "512Mi"
      cpu: "250m"
  
  autoscaling:
    enabled: true
    minReplicas: 2
    maxReplicas: 10
    targetCPUUtilization: 80

# AI Gateway (Python)
aiGateway:
  enabled: true
  replicaCount: 3
  
  image:
    repository: sce-ai-gateway
    tag: latest
  
  resources:
    requests:
      memory: "512Mi"
      cpu: "250m"
    limits:
      memory: "2Gi"
      cpu: "1000m"
  
  autoscaling:
    enabled: true
    minReplicas: 2
    maxReplicas: 10
    targetCPUUtilization: 70
  
  config:
    anthropicApiKey:
      existingSecret: ai-secrets
      key: anthropic-api-key
    openaiApiKey:
      existingSecret: ai-secrets
      key: openai-api-key
    vllmEndpoint: http://vllm-service.omniroute-ai:8000
    rateLimitPerMinute: 100

# vLLM Inference
vllm:
  enabled: true
  
  image:
    repository: vllm/vllm-openai
    tag: latest
  
  model: meta-llama/Llama-3.3-70B-Instruct
  tensorParallelSize: 4
  maxModelLen: 32768
  gpuMemoryUtilization: 0.95
  
  resources:
    limits:
      nvidia.com/gpu: 4
      memory: "320Gi"
      cpu: "32"
  
  persistence:
    enabled: true
    size: 500Gi
    storageClass: premium-rwo

# Temporal
temporal:
  enabled: true
  
  server:
    replicaCount: 3
    image: temporalio/server
    tag: 1.24.0
    
    resources:
      requests:
        memory: "2Gi"
        cpu: "1000m"
      limits:
        memory: "4Gi"
        cpu: "2000m"
    
    config:
      numHistoryShards: 512
  
  web:
    enabled: true
    replicaCount: 1
  
  workers:
    core:
      replicaCount: 5
      taskQueues:
        - omniroute-core
      autoscaling:
        enabled: true
        minReplicas: 3
        maxReplicas: 20
    
    integration:
      replicaCount: 3
      taskQueues:
        - omniroute-integration
      autoscaling:
        enabled: true
        minReplicas: 2
        maxReplicas: 15
    
    ai:
      replicaCount: 2
      taskQueues:
        - omniroute-ai
      autoscaling:
        enabled: true
        minReplicas: 1
        maxReplicas: 10

# n8n
n8n:
  enabled: true
  
  main:
    replicaCount: 2
    image: n8nio/n8n
    tag: latest
    
    resources:
      requests:
        memory: "512Mi"
        cpu: "250m"
      limits:
        memory: "2Gi"
        cpu: "1000m"
    
    env:
      - name: N8N_ENCRYPTION_KEY
        valueFrom:
          secretKeyRef:
            name: n8n-secrets
            key: encryption-key
      - name: EXECUTIONS_MODE
        value: "queue"
  
  worker:
    replicaCount: 5
    autoscaling:
      enabled: true
      minReplicas: 3
      maxReplicas: 15

# Frontend
frontend:
  enabled: true
  replicaCount: 2
  
  image:
    repository: sce-frontend
    tag: latest
  
  resources:
    requests:
      memory: "128Mi"
      cpu: "50m"
    limits:
      memory: "256Mi"
      cpu: "100m"

# Ingress
ingress:
  enabled: true
  className: gce
  
  annotations:
    kubernetes.io/ingress.global-static-ip-name: omniroute-sce-ip
    networking.gke.io/managed-certificates: omniroute-sce-cert
    cloud.google.com/backend-config: '{"default": "sce-backend-config"}'
  
  hosts:
    - host: sce.omniroute.io
      paths:
        - path: /
          pathType: Prefix
          backend:
            service:
              name: frontend
              port: 80
        - path: /api
          pathType: Prefix
          backend:
            service:
              name: api-gateway
              port: 80
```

**Output:**
- Complete Helm chart
- Environment-specific values files
- CI/CD integration scripts
```

---

## Summary

These comprehensive prompts cover the entire implementation of the OmniRoute Service Creation Environment:

1. **Domain Layer**: Aggregates, value objects, repositories following DDD
2. **Temporal Workflows**: Dynamic executor, activities, worker setup
3. **AI Generation**: Prompt-to-service, workflow enhancement
4. **Frontend**: React Flow canvas, node components, state management
5. **Infrastructure**: GKE Terraform, Helm charts, observability

Each prompt includes:
- Clear context and requirements
- Code structure and templates
- Test cases for TDD
- Expected outputs

Use these prompts with Claude, GPT-4, Gemini, or local models (Llama, Mistral) to generate production-ready code that follows XP, DDD, and legacy modernization principles.
