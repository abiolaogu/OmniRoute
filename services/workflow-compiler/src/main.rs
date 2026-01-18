//! Workflow Compiler - Transforms visual workflow DSL into executable Temporal Go code
//! Follows DDD principles with clear domain separation

use axum::{
    extract::State,
    http::StatusCode,
    routing::{get, post},
    Json, Router,
};
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use tokio::net::TcpListener;
use tracing::{info, Level};
use tracing_subscriber::FmtSubscriber;
use uuid::Uuid;

mod compiler;
mod dsl;
mod error;

pub use error::CompilerError;

// =============================================================================
// DOMAIN MODELS
// =============================================================================

/// Workflow definition from visual editor
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WorkflowDefinition {
    pub id: Uuid,
    pub name: String,
    pub version: String,
    pub description: Option<String>,
    pub nodes: Vec<WorkflowNode>,
    pub edges: Vec<WorkflowEdge>,
    pub variables: Vec<Variable>,
    pub triggers: Vec<Trigger>,
}

/// Node in the workflow graph
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WorkflowNode {
    pub id: String,
    pub node_type: NodeType,
    pub label: String,
    pub config: serde_json::Value,
    pub position: Position,
    pub retries: Option<RetryPolicy>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum NodeType {
    Start,
    End,
    Activity,
    Decision,
    ParallelGateway,
    WaitTimer,
    WaitSignal,
    SubWorkflow,
    HttpCall,
    DatabaseQuery,
    Transform,
    Notification,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Position {
    pub x: f64,
    pub y: f64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RetryPolicy {
    pub max_attempts: u32,
    pub initial_interval: String,
    pub max_interval: String,
    pub backoff_coefficient: f64,
}

/// Edge connecting nodes
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WorkflowEdge {
    pub id: String,
    pub source: String,
    pub target: String,
    pub condition: Option<String>,
    pub label: Option<String>,
}

/// Workflow variable
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Variable {
    pub name: String,
    pub var_type: String,
    pub default_value: Option<serde_json::Value>,
}

/// Workflow trigger
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Trigger {
    pub trigger_type: TriggerType,
    pub config: serde_json::Value,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum TriggerType {
    Manual,
    Schedule,
    Webhook,
    Event,
}

// =============================================================================
// COMPILATION OUTPUT
// =============================================================================

/// Compiled workflow output
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CompiledWorkflow {
    pub workflow_code: String,
    pub activity_code: String,
    pub worker_code: String,
    pub test_code: String,
    pub metadata: CompilationMetadata,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CompilationMetadata {
    pub workflow_name: String,
    pub package_name: String,
    pub activities: Vec<String>,
    pub signals: Vec<String>,
    pub queries: Vec<String>,
    pub estimated_complexity: u32,
}

// =============================================================================
// API HANDLERS
// =============================================================================

#[derive(Clone)]
struct AppState {
    compiler: Arc<WorkflowCompiler>,
}

struct WorkflowCompiler {
    templates: handlebars::Handlebars<'static>,
}

impl WorkflowCompiler {
    fn new() -> Self {
        let mut templates = handlebars::Handlebars::new();
        
        // Register templates for Go code generation
        templates.register_template_string("workflow", include_str!("templates/workflow.hbs"))
            .expect("Failed to register workflow template");
        templates.register_template_string("activity", include_str!("templates/activity.hbs"))
            .expect("Failed to register activity template");
        
        Self { templates }
    }
    
    fn compile(&self, definition: &WorkflowDefinition) -> Result<CompiledWorkflow, CompilerError> {
        // Validate workflow
        self.validate(definition)?;
        
        // Optimize graph
        let optimized = self.optimize(definition)?;
        
        // Generate code
        self.generate_code(&optimized)
    }
    
    fn validate(&self, definition: &WorkflowDefinition) -> Result<(), CompilerError> {
        // Check for start and end nodes
        let has_start = definition.nodes.iter().any(|n| matches!(n.node_type, NodeType::Start));
        let has_end = definition.nodes.iter().any(|n| matches!(n.node_type, NodeType::End));
        
        if !has_start {
            return Err(CompilerError::ValidationError("Missing start node".into()));
        }
        if !has_end {
            return Err(CompilerError::ValidationError("Missing end node".into()));
        }
        
        // Check for cycles (simplified)
        // Full implementation would use petgraph for cycle detection
        
        Ok(())
    }
    
    fn optimize(&self, definition: &WorkflowDefinition) -> Result<WorkflowDefinition, CompilerError> {
        // Clone and optimize
        let mut optimized = definition.clone();
        
        // Remove unreachable nodes
        // Merge sequential activities
        // Optimize parallel branches
        
        Ok(optimized)
    }
    
    fn generate_code(&self, definition: &WorkflowDefinition) -> Result<CompiledWorkflow, CompilerError> {
        let package_name = definition.name.to_lowercase().replace(" ", "_");
        
        // Extract activities from nodes
        let activities: Vec<String> = definition.nodes.iter()
            .filter(|n| matches!(n.node_type, NodeType::Activity | NodeType::HttpCall | NodeType::DatabaseQuery))
            .map(|n| format!("{}Activity", to_pascal_case(&n.label)))
            .collect();
        
        // Generate workflow code
        let workflow_code = self.generate_workflow_code(definition, &package_name)?;
        let activity_code = self.generate_activity_code(definition, &package_name)?;
        let worker_code = self.generate_worker_code(definition, &package_name)?;
        let test_code = self.generate_test_code(definition, &package_name)?;
        
        Ok(CompiledWorkflow {
            workflow_code,
            activity_code,
            worker_code,
            test_code,
            metadata: CompilationMetadata {
                workflow_name: definition.name.clone(),
                package_name,
                activities,
                signals: vec![],
                queries: vec![],
                estimated_complexity: definition.nodes.len() as u32,
            },
        })
    }
    
    fn generate_workflow_code(&self, definition: &WorkflowDefinition, package_name: &str) -> Result<String, CompilerError> {
        let workflow_name = to_pascal_case(&definition.name);
        
        Ok(format!(r#"// Generated by OmniRoute Workflow Compiler
// DO NOT EDIT - This file is auto-generated

package {package_name}

import (
    "go.temporal.io/sdk/workflow"
    "time"
)

// {workflow_name}Input defines the workflow input
type {workflow_name}Input struct {{
    // Add input fields based on workflow variables
}}

// {workflow_name}Output defines the workflow output
type {workflow_name}Output struct {{
    Success bool
    Message string
}}

// {workflow_name} is the main workflow function
func {workflow_name}(ctx workflow.Context, input {workflow_name}Input) (*{workflow_name}Output, error) {{
    logger := workflow.GetLogger(ctx)
    logger.Info("{workflow_name} started")
    
    // Activity options
    ao := workflow.ActivityOptions{{
        StartToCloseTimeout: 10 * time.Minute,
    }}
    ctx = workflow.WithActivityOptions(ctx, ao)
    
    // TODO: Generated workflow logic from nodes
    
    return &{workflow_name}Output{{
        Success: true,
        Message: "Workflow completed successfully",
    }}, nil
}}
"#))
    }
    
    fn generate_activity_code(&self, definition: &WorkflowDefinition, package_name: &str) -> Result<String, CompilerError> {
        Ok(format!(r#"// Generated by OmniRoute Workflow Compiler
package {package_name}

import (
    "context"
)

// Activities struct holds activity implementations
type Activities struct {{
    // Add dependencies here
}}

// NewActivities creates a new Activities instance
func NewActivities() *Activities {{
    return &Activities{{}}
}}

// TODO: Generate activity methods from workflow nodes
"#))
    }
    
    fn generate_worker_code(&self, definition: &WorkflowDefinition, package_name: &str) -> Result<String, CompilerError> {
        let workflow_name = to_pascal_case(&definition.name);
        
        Ok(format!(r#"// Generated by OmniRoute Workflow Compiler
package main

import (
    "log"
    "go.temporal.io/sdk/client"
    "go.temporal.io/sdk/worker"
    "{package_name}"
)

func main() {{
    c, err := client.Dial(client.Options{{}})
    if err != nil {{
        log.Fatalln("Unable to create client", err)
    }}
    defer c.Close()

    w := worker.New(c, "{package_name}-task-queue", worker.Options{{}})

    w.RegisterWorkflow({package_name}.{workflow_name})
    
    activities := {package_name}.NewActivities()
    w.RegisterActivity(activities)

    err = w.Run(worker.InterruptCh())
    if err != nil {{
        log.Fatalln("Unable to start worker", err)
    }}
}}
"#))
    }
    
    fn generate_test_code(&self, definition: &WorkflowDefinition, package_name: &str) -> Result<String, CompilerError> {
        let workflow_name = to_pascal_case(&definition.name);
        
        Ok(format!(r#"// Generated by OmniRoute Workflow Compiler
package {package_name}

import (
    "testing"
    "github.com/stretchr/testify/require"
    "go.temporal.io/sdk/testsuite"
)

func Test{workflow_name}(t *testing.T) {{
    testSuite := &testsuite.WorkflowTestSuite{{}}
    env := testSuite.NewTestWorkflowEnvironment()

    env.RegisterWorkflow({workflow_name})
    activities := NewActivities()
    env.RegisterActivity(activities)

    env.ExecuteWorkflow({workflow_name}, {workflow_name}Input{{}})

    require.True(t, env.IsWorkflowCompleted())
    require.NoError(t, env.GetWorkflowError())
}}
"#))
    }
}

fn to_pascal_case(s: &str) -> String {
    s.split(|c: char| c.is_whitespace() || c == '_' || c == '-')
        .map(|word| {
            let mut chars = word.chars();
            match chars.next() {
                None => String::new(),
                Some(c) => c.to_uppercase().chain(chars).collect(),
            }
        })
        .collect()
}

// API Handlers

#[derive(Deserialize)]
struct CompileRequest {
    workflow: WorkflowDefinition,
}

#[derive(Serialize)]
struct CompileResponse {
    success: bool,
    compiled: Option<CompiledWorkflow>,
    error: Option<String>,
}

async fn health() -> &'static str {
    "OK"
}

async fn compile_workflow(
    State(state): State<AppState>,
    Json(request): Json<CompileRequest>,
) -> Result<Json<CompileResponse>, StatusCode> {
    match state.compiler.compile(&request.workflow) {
        Ok(compiled) => Ok(Json(CompileResponse {
            success: true,
            compiled: Some(compiled),
            error: None,
        })),
        Err(e) => Ok(Json(CompileResponse {
            success: false,
            compiled: None,
            error: Some(e.to_string()),
        })),
    }
}

async fn validate_workflow(
    State(state): State<AppState>,
    Json(request): Json<CompileRequest>,
) -> Result<Json<serde_json::Value>, StatusCode> {
    match state.compiler.validate(&request.workflow) {
        Ok(_) => Ok(Json(serde_json::json!({
            "valid": true,
            "errors": []
        }))),
        Err(e) => Ok(Json(serde_json::json!({
            "valid": false,
            "errors": [e.to_string()]
        }))),
    }
}

#[tokio::main]
async fn main() {
    // Initialize tracing
    let subscriber = FmtSubscriber::builder()
        .with_max_level(Level::INFO)
        .finish();
    tracing::subscriber::set_global_default(subscriber).expect("setting default subscriber failed");
    
    let state = AppState {
        compiler: Arc::new(WorkflowCompiler::new()),
    };
    
    let app = Router::new()
        .route("/health", get(health))
        .route("/api/v1/compile", post(compile_workflow))
        .route("/api/v1/validate", post(validate_workflow))
        .with_state(state);
    
    let port = std::env::var("PORT").unwrap_or_else(|_| "8130".to_string());
    let listener = TcpListener::bind(format!("0.0.0.0:{}", port)).await.unwrap();
    
    info!("Workflow Compiler listening on port {}", port);
    axum::serve(listener, app).await.unwrap();
}
