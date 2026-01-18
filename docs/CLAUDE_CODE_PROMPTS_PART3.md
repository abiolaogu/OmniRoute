# OmniRoute SCE - Claude Code Implementation Prompts
## Part 3: Workflow Compiler (Rust), AI Gateway (Python), and Frontend (React)

---

# PHASE 3: WORKFLOW COMPILER (RUST)

## Prompt 3.1: Rust Workflow Compiler Core

```
Implement the Workflow Compiler service in Rust that transforms visual workflow DSL into executable Temporal Go code.

## Project Structure
```
services/workflow-compiler/
├── Cargo.toml
├── src/
│   ├── main.rs
│   ├── lib.rs
│   ├── api/
│   │   ├── mod.rs
│   │   ├── server.rs
│   │   └── handlers.rs
│   ├── compiler/
│   │   ├── mod.rs
│   │   ├── parser.rs
│   │   ├── validator.rs
│   │   ├── optimizer.rs
│   │   └── codegen.rs
│   ├── dsl/
│   │   ├── mod.rs
│   │   ├── types.rs
│   │   └── graph.rs
│   ├── templates/
│   │   ├── mod.rs
│   │   ├── workflow.hbs
│   │   ├── activity.hbs
│   │   └── node_handlers.hbs
│   └── error.rs
├── tests/
│   ├── compiler_tests.rs
│   └── fixtures/
└── benches/
    └── compiler_bench.rs
```

## Implementation

### 1. Cargo.toml

```toml
[package]
name = "workflow-compiler"
version = "0.1.0"
edition = "2021"
authors = ["OmniRoute Team"]

[lib]
name = "workflow_compiler"
path = "src/lib.rs"

[[bin]]
name = "workflow-compiler-server"
path = "src/main.rs"

[dependencies]
# Async runtime
tokio = { version = "1.35", features = ["full"] }

# Serialization
serde = { version = "1.0", features = ["derive"] }
serde_json = "1.0"

# Error handling
thiserror = "1.0"
anyhow = "1.0"

# Graph algorithms
petgraph = "0.6"

# Template engine
handlebars = "5.0"

# HTTP server
axum = "0.7"
tower = "0.4"
tower-http = { version = "0.5", features = ["cors", "trace"] }

# Tracing
tracing = "0.1"
tracing-subscriber = { version = "0.3", features = ["env-filter", "json"] }
tracing-opentelemetry = "0.22"
opentelemetry = "0.21"
opentelemetry-otlp = "0.14"

# Validation
validator = { version = "0.16", features = ["derive"] }

# Utilities
uuid = { version = "1.6", features = ["v4", "serde"] }
chrono = { version = "0.4", features = ["serde"] }

[dev-dependencies]
criterion = "0.5"
proptest = "1.4"
insta = { version = "1.34", features = ["json"] }
tokio-test = "0.4"

[[bench]]
name = "compiler_bench"
harness = false
```

### 2. DSL Types (src/dsl/types.rs)

```rust
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use uuid::Uuid;
use validator::Validate;

/// The complete workflow DSL representation
#[derive(Debug, Clone, Serialize, Deserialize, Validate)]
pub struct WorkflowDSL {
    pub id: Uuid,
    pub name: String,
    pub tenant_id: Uuid,
    pub version: u32,
    
    #[validate(length(min = 1))]
    pub nodes: Vec<WorkflowNode>,
    pub edges: Vec<WorkflowEdge>,
    pub variables: Vec<WorkflowVariable>,
    pub triggers: Vec<WorkflowTrigger>,
    pub error_policy: ErrorHandlingPolicy,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WorkflowNode {
    pub id: String,
    pub node_type: NodeType,
    pub label: String,
    pub position: Position,
    pub config: NodeConfig,
    pub retry_policy: Option<RetryPolicy>,
    pub timeout_seconds: Option<u64>,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, Eq)]
#[serde(rename_all = "snake_case")]
pub enum NodeType {
    Activity,
    Subflow,
    AiAction,
    N8n,
    Decision,
    Parallel,
    Wait,
    HumanTask,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(tag = "type", rename_all = "snake_case")]
pub enum NodeConfig {
    Activity(ActivityConfig),
    Subflow(SubflowConfig),
    AiAction(AiActionConfig),
    N8n(N8nConfig),
    Decision(DecisionConfig),
    Parallel(ParallelConfig),
    Wait(WaitConfig),
    HumanTask(HumanTaskConfig),
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ActivityConfig {
    pub activity_name: String,
    pub input_mapping: HashMap<String, String>,
    pub output_variable: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AiActionConfig {
    pub provider: AiProvider,
    pub model: String,
    pub prompt_template: String,
    pub temperature: f64,
    pub max_tokens: u32,
    pub use_local_model: bool,
    pub output_schema: Option<serde_json::Value>,
    pub output_variable: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "lowercase")]
pub enum AiProvider {
    Anthropic,
    OpenAI,
    Google,
    Local,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct N8nConfig {
    pub workflow_id: String,
    pub webhook_path: Option<String>,
    pub input_mapping: HashMap<String, String>,
    pub wait_for_completion: bool,
    pub timeout_seconds: Option<u64>,
    pub output_variable: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DecisionConfig {
    pub condition_expression: String,
    pub branches: Vec<DecisionBranch>,
    pub default_branch: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DecisionBranch {
    pub id: String,
    pub condition: String,
    pub label: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ParallelConfig {
    pub branches: Vec<String>,  // Node IDs to execute in parallel
    pub wait_for_all: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WaitConfig {
    pub wait_type: WaitType,
    pub duration_seconds: Option<u64>,
    pub signal_name: Option<String>,
    pub timeout_seconds: Option<u64>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum WaitType {
    Timer,
    Signal,
    Condition,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct HumanTaskConfig {
    pub task_type: String,
    pub title: String,
    pub description: Option<String>,
    pub assignee_expression: Option<String>,
    pub form_schema: Option<serde_json::Value>,
    pub timeout_seconds: Option<u64>,
    pub output_variable: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SubflowConfig {
    pub workflow_id: String,
    pub input_mapping: HashMap<String, String>,
    pub output_variable: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WorkflowEdge {
    pub id: String,
    pub source: String,
    pub target: String,
    pub source_handle: Option<String>,
    pub target_handle: Option<String>,
    pub condition: Option<EdgeCondition>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct EdgeCondition {
    pub expression: String,
    pub label: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WorkflowVariable {
    pub name: String,
    pub variable_type: VariableType,
    pub default_value: Option<serde_json::Value>,
    pub required: bool,
    pub description: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "lowercase")]
pub enum VariableType {
    String,
    Number,
    Boolean,
    Object,
    Array,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WorkflowTrigger {
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

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RetryPolicy {
    pub max_attempts: u32,
    pub initial_interval_ms: u64,
    pub backoff_coefficient: f64,
    pub max_interval_ms: u64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ErrorHandlingPolicy {
    pub on_error: ErrorAction,
    pub compensation_workflow_id: Option<String>,
    pub notification_channel: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "lowercase")]
pub enum ErrorAction {
    Fail,
    Compensate,
    Ignore,
    Retry,
}

#[derive(Debug, Clone, Copy, Serialize, Deserialize)]
pub struct Position {
    pub x: f64,
    pub y: f64,
}
```

### 3. Graph Module (src/dsl/graph.rs)

```rust
use petgraph::graph::{DiGraph, NodeIndex};
use petgraph::algo::{toposort, is_cyclic_directed};
use petgraph::visit::EdgeRef;
use std::collections::HashMap;

use crate::dsl::types::{WorkflowDSL, WorkflowNode, WorkflowEdge, NodeType};
use crate::error::{CompilerError, CompilerResult};

/// Represents the workflow as a directed graph
pub struct WorkflowGraph {
    graph: DiGraph<String, EdgeData>,
    node_map: HashMap<String, NodeIndex>,
    node_data: HashMap<String, WorkflowNode>,
}

struct EdgeData {
    edge_id: String,
    condition: Option<String>,
}

impl WorkflowGraph {
    /// Build a graph from the DSL
    pub fn from_dsl(dsl: &WorkflowDSL) -> CompilerResult<Self> {
        let mut graph = DiGraph::new();
        let mut node_map = HashMap::new();
        let mut node_data = HashMap::new();
        
        // Add nodes
        for node in &dsl.nodes {
            let idx = graph.add_node(node.id.clone());
            node_map.insert(node.id.clone(), idx);
            node_data.insert(node.id.clone(), node.clone());
        }
        
        // Add edges
        for edge in &dsl.edges {
            let source_idx = node_map.get(&edge.source)
                .ok_or_else(|| CompilerError::InvalidEdge {
                    edge_id: edge.id.clone(),
                    reason: format!("Source node '{}' not found", edge.source),
                })?;
            
            let target_idx = node_map.get(&edge.target)
                .ok_or_else(|| CompilerError::InvalidEdge {
                    edge_id: edge.id.clone(),
                    reason: format!("Target node '{}' not found", edge.target),
                })?;
            
            graph.add_edge(
                *source_idx,
                *target_idx,
                EdgeData {
                    edge_id: edge.id.clone(),
                    condition: edge.condition.as_ref().map(|c| c.expression.clone()),
                },
            );
        }
        
        Ok(Self {
            graph,
            node_map,
            node_data,
        })
    }
    
    /// Check if the graph has cycles
    pub fn has_cycle(&self) -> bool {
        is_cyclic_directed(&self.graph)
    }
    
    /// Get topological order of nodes
    pub fn topological_order(&self) -> CompilerResult<Vec<String>> {
        match toposort(&self.graph, None) {
            Ok(order) => {
                Ok(order.into_iter()
                    .map(|idx| self.graph[idx].clone())
                    .collect())
            }
            Err(_) => Err(CompilerError::CyclicGraph),
        }
    }
    
    /// Get entry nodes (nodes with no incoming edges)
    pub fn entry_nodes(&self) -> Vec<String> {
        self.node_map.iter()
            .filter(|(_, idx)| {
                self.graph.edges_directed(**idx, petgraph::Direction::Incoming).count() == 0
            })
            .map(|(id, _)| id.clone())
            .collect()
    }
    
    /// Get outgoing edges for a node
    pub fn outgoing_edges(&self, node_id: &str) -> Vec<(String, Option<String>)> {
        let Some(idx) = self.node_map.get(node_id) else {
            return vec![];
        };
        
        self.graph.edges(*idx)
            .map(|edge| {
                let target_idx = edge.target();
                let target_id = self.graph[target_idx].clone();
                let condition = edge.weight().condition.clone();
                (target_id, condition)
            })
            .collect()
    }
    
    /// Get node data by ID
    pub fn get_node(&self, node_id: &str) -> Option<&WorkflowNode> {
        self.node_data.get(node_id)
    }
    
    /// Find parallel branches starting from a parallel node
    pub fn find_parallel_branches(&self, parallel_node_id: &str) -> Vec<Vec<String>> {
        let mut branches = Vec::new();
        
        for (target_id, _) in self.outgoing_edges(parallel_node_id) {
            let mut branch = vec![target_id.clone()];
            let mut current = target_id;
            
            // Follow the branch until we hit a join or end
            loop {
                let outgoing = self.outgoing_edges(&current);
                if outgoing.len() != 1 {
                    break;
                }
                let (next_id, _) = &outgoing[0];
                branch.push(next_id.clone());
                current = next_id.clone();
            }
            
            branches.push(branch);
        }
        
        branches
    }
    
    /// Validate graph structure
    pub fn validate(&self) -> CompilerResult<()> {
        // Check for cycles
        if self.has_cycle() {
            return Err(CompilerError::CyclicGraph);
        }
        
        // Check for orphan nodes (disconnected from the flow)
        let entry_nodes = self.entry_nodes();
        if entry_nodes.is_empty() {
            return Err(CompilerError::NoEntryNode);
        }
        
        // Validate decision nodes have multiple outputs
        for (node_id, node) in &self.node_data {
            if matches!(node.node_type, NodeType::Decision) {
                let outgoing = self.outgoing_edges(node_id);
                if outgoing.len() < 2 {
                    return Err(CompilerError::InvalidNode {
                        node_id: node_id.clone(),
                        reason: "Decision node must have at least 2 outgoing edges".to_string(),
                    });
                }
            }
        }
        
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_cycle_detection() {
        // Create DSL with cycle
        let dsl = create_cyclic_dsl();
        let graph = WorkflowGraph::from_dsl(&dsl).unwrap();
        assert!(graph.has_cycle());
    }
    
    #[test]
    fn test_topological_order() {
        // Create valid linear DSL
        let dsl = create_linear_dsl();
        let graph = WorkflowGraph::from_dsl(&dsl).unwrap();
        let order = graph.topological_order().unwrap();
        assert_eq!(order, vec!["node1", "node2", "node3"]);
    }
}
```

### 4. Code Generator (src/compiler/codegen.rs)

```rust
use handlebars::Handlebars;
use serde::Serialize;
use std::collections::HashMap;

use crate::dsl::types::*;
use crate::dsl::graph::WorkflowGraph;
use crate::error::{CompilerError, CompilerResult};

/// Generated workflow code
pub struct GeneratedWorkflow {
    pub code: String,
    pub activity_stubs: Vec<String>,
    pub required_workers: Vec<String>,
}

/// Code generator for Temporal workflows
pub struct CodeGenerator<'a> {
    handlebars: Handlebars<'a>,
}

impl<'a> CodeGenerator<'a> {
    pub fn new() -> CompilerResult<Self> {
        let mut handlebars = Handlebars::new();
        
        // Register templates
        handlebars.register_template_string("workflow", include_str!("../templates/workflow.hbs"))
            .map_err(|e| CompilerError::TemplateError(e.to_string()))?;
        handlebars.register_template_string("activity_call", include_str!("../templates/activity.hbs"))
            .map_err(|e| CompilerError::TemplateError(e.to_string()))?;
        handlebars.register_template_string("ai_action", include_str!("../templates/ai_action.hbs"))
            .map_err(|e| CompilerError::TemplateError(e.to_string()))?;
        handlebars.register_template_string("n8n_call", include_str!("../templates/n8n.hbs"))
            .map_err(|e| CompilerError::TemplateError(e.to_string()))?;
        handlebars.register_template_string("decision", include_str!("../templates/decision.hbs"))
            .map_err(|e| CompilerError::TemplateError(e.to_string()))?;
        handlebars.register_template_string("parallel", include_str!("../templates/parallel.hbs"))
            .map_err(|e| CompilerError::TemplateError(e.to_string()))?;
        handlebars.register_template_string("wait", include_str!("../templates/wait.hbs"))
            .map_err(|e| CompilerError::TemplateError(e.to_string()))?;
        handlebars.register_template_string("human_task", include_str!("../templates/human_task.hbs"))
            .map_err(|e| CompilerError::TemplateError(e.to_string()))?;
        
        // Register helpers
        handlebars.register_helper("snake_case", Box::new(snake_case_helper));
        handlebars.register_helper("pascal_case", Box::new(pascal_case_helper));
        
        Ok(Self { handlebars })
    }
    
    pub fn generate(&self, dsl: &WorkflowDSL) -> CompilerResult<GeneratedWorkflow> {
        let graph = WorkflowGraph::from_dsl(dsl)?;
        graph.validate()?;
        
        let execution_order = graph.topological_order()?;
        
        let mut node_code_blocks = Vec::new();
        let mut activity_stubs = Vec::new();
        let mut required_workers = vec!["omniroute-core".to_string()];
        
        for node_id in &execution_order {
            let node = graph.get_node(node_id)
                .ok_or_else(|| CompilerError::InvalidNode {
                    node_id: node_id.clone(),
                    reason: "Node not found in graph".to_string(),
                })?;
            
            let (code, stubs, workers) = self.generate_node_code(node, &graph)?;
            node_code_blocks.push(code);
            activity_stubs.extend(stubs);
            required_workers.extend(workers);
        }
        
        // Generate the main workflow
        let workflow_context = WorkflowContext {
            name: sanitize_name(&dsl.name),
            workflow_id: dsl.id.to_string(),
            tenant_id: dsl.tenant_id.to_string(),
            variables: dsl.variables.iter().map(|v| VariableContext {
                name: v.name.clone(),
                go_type: variable_type_to_go(&v.variable_type),
                default_value: v.default_value.as_ref().map(|v| format!("{}", v)),
                required: v.required,
            }).collect(),
            node_code: node_code_blocks.join("\n\n"),
            error_policy: match dsl.error_policy.on_error {
                ErrorAction::Fail => "fail",
                ErrorAction::Compensate => "compensate",
                ErrorAction::Ignore => "ignore",
                ErrorAction::Retry => "retry",
            }.to_string(),
            compensation_workflow: dsl.error_policy.compensation_workflow_id.clone(),
        };
        
        let code = self.handlebars.render("workflow", &workflow_context)
            .map_err(|e| CompilerError::CodeGenError(e.to_string()))?;
        
        // Deduplicate
        activity_stubs.sort();
        activity_stubs.dedup();
        required_workers.sort();
        required_workers.dedup();
        
        Ok(GeneratedWorkflow {
            code,
            activity_stubs,
            required_workers,
        })
    }
    
    fn generate_node_code(
        &self,
        node: &WorkflowNode,
        graph: &WorkflowGraph,
    ) -> CompilerResult<(String, Vec<String>, Vec<String>)> {
        let mut activity_stubs = Vec::new();
        let mut required_workers = Vec::new();
        
        let code = match &node.config {
            NodeConfig::Activity(config) => {
                activity_stubs.push(config.activity_name.clone());
                
                let ctx = ActivityContext {
                    node_id: node.id.clone(),
                    activity_name: config.activity_name.clone(),
                    input_mapping: format_input_mapping(&config.input_mapping),
                    output_variable: config.output_variable.clone().unwrap_or_else(|| format!("result_{}", sanitize_name(&node.id))),
                    retry_policy: node.retry_policy.as_ref().map(format_retry_policy),
                    timeout: node.timeout_seconds,
                };
                
                self.handlebars.render("activity_call", &ctx)
                    .map_err(|e| CompilerError::CodeGenError(e.to_string()))?
            }
            
            NodeConfig::AiAction(config) => {
                required_workers.push("omniroute-ai".to_string());
                
                let ctx = AiActionContext {
                    node_id: node.id.clone(),
                    provider: format!("{:?}", config.provider).to_lowercase(),
                    model: config.model.clone(),
                    prompt_template: escape_go_string(&config.prompt_template),
                    temperature: config.temperature,
                    max_tokens: config.max_tokens,
                    use_local_model: config.use_local_model,
                    output_variable: config.output_variable.clone().unwrap_or_else(|| format!("aiResult_{}", sanitize_name(&node.id))),
                    timeout: node.timeout_seconds.unwrap_or(120),
                };
                
                self.handlebars.render("ai_action", &ctx)
                    .map_err(|e| CompilerError::CodeGenError(e.to_string()))?
            }
            
            NodeConfig::N8n(config) => {
                required_workers.push("omniroute-integration".to_string());
                
                let ctx = N8nContext {
                    node_id: node.id.clone(),
                    workflow_id: config.workflow_id.clone(),
                    webhook_path: config.webhook_path.clone(),
                    input_mapping: format_input_mapping(&config.input_mapping),
                    wait_for_completion: config.wait_for_completion,
                    output_variable: config.output_variable.clone().unwrap_or_else(|| format!("n8nResult_{}", sanitize_name(&node.id))),
                    timeout: config.timeout_seconds.unwrap_or(300),
                };
                
                self.handlebars.render("n8n_call", &ctx)
                    .map_err(|e| CompilerError::CodeGenError(e.to_string()))?
            }
            
            NodeConfig::Decision(config) => {
                let outgoing = graph.outgoing_edges(&node.id);
                let branches: Vec<DecisionBranchContext> = config.branches.iter()
                    .map(|b| {
                        let target = outgoing.iter()
                            .find(|(_, cond)| cond.as_ref().map(|c| c == &b.condition).unwrap_or(false))
                            .map(|(t, _)| t.clone());
                        
                        DecisionBranchContext {
                            condition: b.condition.clone(),
                            label: b.label.clone(),
                            target_node: target,
                        }
                    })
                    .collect();
                
                let ctx = DecisionContext {
                    node_id: node.id.clone(),
                    condition_expression: config.condition_expression.clone(),
                    branches,
                    default_branch: config.default_branch.clone(),
                };
                
                self.handlebars.render("decision", &ctx)
                    .map_err(|e| CompilerError::CodeGenError(e.to_string()))?
            }
            
            NodeConfig::Parallel(config) => {
                let branches = graph.find_parallel_branches(&node.id);
                
                let ctx = ParallelContext {
                    node_id: node.id.clone(),
                    branches: branches.into_iter().enumerate()
                        .map(|(i, nodes)| ParallelBranchContext {
                            index: i,
                            nodes,
                        })
                        .collect(),
                    wait_for_all: config.wait_for_all,
                };
                
                self.handlebars.render("parallel", &ctx)
                    .map_err(|e| CompilerError::CodeGenError(e.to_string()))?
            }
            
            NodeConfig::Wait(config) => {
                let ctx = WaitContext {
                    node_id: node.id.clone(),
                    wait_type: format!("{:?}", config.wait_type).to_lowercase(),
                    duration_seconds: config.duration_seconds,
                    signal_name: config.signal_name.clone(),
                    timeout_seconds: config.timeout_seconds,
                };
                
                self.handlebars.render("wait", &ctx)
                    .map_err(|e| CompilerError::CodeGenError(e.to_string()))?
            }
            
            NodeConfig::HumanTask(config) => {
                let ctx = HumanTaskContext {
                    node_id: node.id.clone(),
                    task_type: config.task_type.clone(),
                    title: config.title.clone(),
                    description: config.description.clone(),
                    assignee_expression: config.assignee_expression.clone(),
                    timeout_seconds: config.timeout_seconds.unwrap_or(86400), // 24h default
                    output_variable: config.output_variable.clone().unwrap_or_else(|| format!("humanResult_{}", sanitize_name(&node.id))),
                };
                
                self.handlebars.render("human_task", &ctx)
                    .map_err(|e| CompilerError::CodeGenError(e.to_string()))?
            }
            
            NodeConfig::Subflow(config) => {
                let ctx = SubflowContext {
                    node_id: node.id.clone(),
                    workflow_id: config.workflow_id.clone(),
                    input_mapping: format_input_mapping(&config.input_mapping),
                    output_variable: config.output_variable.clone().unwrap_or_else(|| format!("subflowResult_{}", sanitize_name(&node.id))),
                };
                
                self.handlebars.render("subflow", &ctx)
                    .map_err(|e| CompilerError::CodeGenError(e.to_string()))?
            }
        };
        
        Ok((code, activity_stubs, required_workers))
    }
}

// Template context structs
#[derive(Serialize)]
struct WorkflowContext {
    name: String,
    workflow_id: String,
    tenant_id: String,
    variables: Vec<VariableContext>,
    node_code: String,
    error_policy: String,
    compensation_workflow: Option<String>,
}

#[derive(Serialize)]
struct VariableContext {
    name: String,
    go_type: String,
    default_value: Option<String>,
    required: bool,
}

#[derive(Serialize)]
struct ActivityContext {
    node_id: String,
    activity_name: String,
    input_mapping: String,
    output_variable: String,
    retry_policy: Option<String>,
    timeout: Option<u64>,
}

#[derive(Serialize)]
struct AiActionContext {
    node_id: String,
    provider: String,
    model: String,
    prompt_template: String,
    temperature: f64,
    max_tokens: u32,
    use_local_model: bool,
    output_variable: String,
    timeout: u64,
}

// ... more context structs

// Helper functions
fn sanitize_name(name: &str) -> String {
    name.chars()
        .map(|c| if c.is_alphanumeric() { c } else { '_' })
        .collect()
}

fn variable_type_to_go(vt: &VariableType) -> String {
    match vt {
        VariableType::String => "string".to_string(),
        VariableType::Number => "float64".to_string(),
        VariableType::Boolean => "bool".to_string(),
        VariableType::Object => "map[string]interface{}".to_string(),
        VariableType::Array => "[]interface{}".to_string(),
    }
}

fn escape_go_string(s: &str) -> String {
    s.replace('\\', "\\\\")
        .replace('"', "\\\"")
        .replace('\n', "\\n")
        .replace('\t', "\\t")
}

fn format_input_mapping(mapping: &HashMap<String, String>) -> String {
    if mapping.is_empty() {
        return "nil".to_string();
    }
    
    let entries: Vec<String> = mapping.iter()
        .map(|(k, v)| format!("\"{}\": {}", k, v))
        .collect();
    
    format!("map[string]interface{{}}{{{}}}", entries.join(", "))
}

fn format_retry_policy(policy: &RetryPolicy) -> String {
    format!(
        r#"&temporal.RetryPolicy{{
            MaximumAttempts: {},
            InitialInterval: {} * time.Millisecond,
            BackoffCoefficient: {},
            MaximumInterval: {} * time.Millisecond,
        }}"#,
        policy.max_attempts,
        policy.initial_interval_ms,
        policy.backoff_coefficient,
        policy.max_interval_ms,
    )
}
```

### 5. Workflow Template (src/templates/workflow.hbs)

```handlebars
// Code generated by OmniRoute SCE Workflow Compiler. DO NOT EDIT.
// Workflow: {{name}}
// ID: {{workflow_id}}
// Generated at: {{now}}

package workflows

import (
    "context"
    "fmt"
    "time"
    
    "go.temporal.io/sdk/workflow"
    "go.temporal.io/sdk/temporal"
    
    "github.com/omniroute/sce/services/workflow-executor/internal/activities"
    "github.com/omniroute/sce/services/workflow-executor/internal/ai"
    "github.com/omniroute/sce/services/workflow-executor/internal/integration"
)

// {{pascal_case name}}Input is the input for the workflow
type {{pascal_case name}}Input struct {
    TenantID      string                 `json:"tenant_id"`
    CorrelationID string                 `json:"correlation_id"`
{{#each variables}}
{{#if required}}
    {{pascal_case name}} {{go_type}} `json:"{{name}}"`
{{else}}
    {{pascal_case name}} *{{go_type}} `json:"{{name}},omitempty"`
{{/if}}
{{/each}}
    Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// {{pascal_case name}}Output is the output of the workflow
type {{pascal_case name}}Output struct {
    Status       string                 `json:"status"`
    Results      map[string]interface{} `json:"results"`
    CompletedAt  time.Time              `json:"completed_at"`
    Error        *string                `json:"error,omitempty"`
}

// {{pascal_case name}}Workflow is a user-defined service workflow
func {{pascal_case name}}Workflow(ctx workflow.Context, input *{{pascal_case name}}Input) (*{{pascal_case name}}Output, error) {
    logger := workflow.GetLogger(ctx)
    logger.Info("Starting workflow", "tenant_id", input.TenantID, "correlation_id", input.CorrelationID)
    
    // Initialize execution context
    execCtx := &ExecutionContext{
        TenantID:      input.TenantID,
        CorrelationID: input.CorrelationID,
        Variables:     make(map[string]interface{}),
        Results:       make(map[string]interface{}),
    }
    
    // Copy input variables
{{#each variables}}
{{#if required}}
    execCtx.Variables["{{name}}"] = input.{{pascal_case name}}
{{else}}
    if input.{{pascal_case name}} != nil {
        execCtx.Variables["{{name}}"] = *input.{{pascal_case name}}
    }{{#if default_value}} else {
        execCtx.Variables["{{name}}"] = {{default_value}}
    }{{/if}}
{{/if}}
{{/each}}

    // Activity options
    defaultActivityOptions := workflow.ActivityOptions{
        StartToCloseTimeout: 5 * time.Minute,
        RetryPolicy: &temporal.RetryPolicy{
            MaximumAttempts: 3,
            InitialInterval: 1 * time.Second,
            BackoffCoefficient: 2.0,
            MaximumInterval: 30 * time.Second,
        },
    }
    ctx = workflow.WithActivityOptions(ctx, defaultActivityOptions)

    var err error
    
    // === Generated Node Code ===
    
{{{node_code}}}
    
    // === End Generated Node Code ===
    
    return &{{pascal_case name}}Output{
        Status:      "completed",
        Results:     execCtx.Results,
        CompletedAt: workflow.Now(ctx),
    }, nil
}

// ExecutionContext holds workflow execution state
type ExecutionContext struct {
    TenantID      string
    CorrelationID string
    Variables     map[string]interface{}
    Results       map[string]interface{}
}

{{#if (eq error_policy "compensate")}}
// compensate runs the compensation workflow on error
func compensate(ctx workflow.Context, execCtx *ExecutionContext, originalErr error) (*{{pascal_case name}}Output, error) {
    logger := workflow.GetLogger(ctx)
    logger.Error("Workflow failed, starting compensation", "error", originalErr)
    
    {{#if compensation_workflow}}
    // Execute compensation workflow
    childCtx := workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
        WorkflowID: fmt.Sprintf("compensation-%s", execCtx.CorrelationID),
    })
    
    var compensationResult interface{}
    err := workflow.ExecuteChildWorkflow(childCtx, "{{compensation_workflow}}", execCtx).Get(ctx, &compensationResult)
    if err != nil {
        logger.Error("Compensation workflow failed", "error", err)
    }
    {{/if}}
    
    errStr := originalErr.Error()
    return &{{pascal_case name}}Output{
        Status:      "failed",
        Results:     execCtx.Results,
        CompletedAt: workflow.Now(ctx),
        Error:       &errStr,
    }, originalErr
}
{{/if}}
```

### 6. HTTP Server (src/api/server.rs)

```rust
use axum::{
    routing::{get, post},
    Router,
    Json,
    extract::State,
    http::StatusCode,
};
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use tower_http::trace::TraceLayer;
use tracing::info;

use crate::compiler::Compiler;
use crate::dsl::types::WorkflowDSL;
use crate::error::CompilerError;

pub struct AppState {
    compiler: Compiler,
}

#[derive(Deserialize)]
pub struct CompileRequest {
    workflow: WorkflowDSL,
}

#[derive(Serialize)]
pub struct CompileResponse {
    code: String,
    activity_stubs: Vec<String>,
    required_workers: Vec<String>,
}

#[derive(Serialize)]
pub struct ErrorResponse {
    error: String,
    message: String,
    details: Option<serde_json::Value>,
}

pub async fn run_server(addr: &str) -> Result<(), Box<dyn std::error::Error>> {
    let compiler = Compiler::new()?;
    let state = Arc::new(AppState { compiler });
    
    let app = Router::new()
        .route("/health", get(health_check))
        .route("/api/v1/compile", post(compile_workflow))
        .route("/api/v1/validate", post(validate_workflow))
        .with_state(state)
        .layer(TraceLayer::new_for_http());
    
    info!("Starting workflow compiler server on {}", addr);
    
    let listener = tokio::net::TcpListener::bind(addr).await?;
    axum::serve(listener, app).await?;
    
    Ok(())
}

async fn health_check() -> &'static str {
    "OK"
}

async fn compile_workflow(
    State(state): State<Arc<AppState>>,
    Json(request): Json<CompileRequest>,
) -> Result<Json<CompileResponse>, (StatusCode, Json<ErrorResponse>)> {
    match state.compiler.compile(&request.workflow) {
        Ok(result) => Ok(Json(CompileResponse {
            code: result.code,
            activity_stubs: result.activity_stubs,
            required_workers: result.required_workers,
        })),
        Err(e) => {
            let (status, error_response) = map_error(e);
            Err((status, Json(error_response)))
        }
    }
}

async fn validate_workflow(
    State(state): State<Arc<AppState>>,
    Json(request): Json<CompileRequest>,
) -> Result<Json<serde_json::Value>, (StatusCode, Json<ErrorResponse>)> {
    match state.compiler.validate(&request.workflow) {
        Ok(warnings) => Ok(Json(serde_json::json!({
            "valid": true,
            "warnings": warnings,
        }))),
        Err(e) => {
            let (status, error_response) = map_error(e);
            Err((status, Json(error_response)))
        }
    }
}

fn map_error(e: CompilerError) -> (StatusCode, ErrorResponse) {
    match e {
        CompilerError::CyclicGraph => (
            StatusCode::BAD_REQUEST,
            ErrorResponse {
                error: "cyclic_graph".to_string(),
                message: "Workflow contains a cycle".to_string(),
                details: None,
            },
        ),
        CompilerError::InvalidNode { node_id, reason } => (
            StatusCode::BAD_REQUEST,
            ErrorResponse {
                error: "invalid_node".to_string(),
                message: format!("Invalid node '{}': {}", node_id, reason),
                details: Some(serde_json::json!({ "node_id": node_id })),
            },
        ),
        CompilerError::InvalidEdge { edge_id, reason } => (
            StatusCode::BAD_REQUEST,
            ErrorResponse {
                error: "invalid_edge".to_string(),
                message: format!("Invalid edge '{}': {}", edge_id, reason),
                details: Some(serde_json::json!({ "edge_id": edge_id })),
            },
        ),
        _ => (
            StatusCode::INTERNAL_SERVER_ERROR,
            ErrorResponse {
                error: "internal_error".to_string(),
                message: "An internal error occurred".to_string(),
                details: None,
            },
        ),
    }
}
```

## Deliverables

1. Complete Rust project with all modules
2. DSL type definitions with validation
3. Graph algorithms (cycle detection, topological sort)
4. Code generator with Handlebars templates
5. HTTP server with Axum
6. Comprehensive tests with fixtures
7. Benchmarks for performance testing
```

---

# PHASE 4: AI GATEWAY (PYTHON)

## Prompt 4.1: AI Gateway Service

```
Implement the AI Gateway service in Python using FastAPI.

## Project Structure
```
services/ai-gateway/
├── pyproject.toml
├── app/
│   ├── __init__.py
│   ├── main.py
│   ├── config.py
│   ├── models/
│   │   ├── __init__.py
│   │   ├── requests.py
│   │   └── responses.py
│   ├── providers/
│   │   ├── __init__.py
│   │   ├── base.py
│   │   ├── anthropic.py
│   │   ├── openai.py
│   │   └── vllm.py
│   ├── services/
│   │   ├── __init__.py
│   │   ├── generator.py
│   │   ├── enhancer.py
│   │   └── router.py
│   ├── prompts/
│   │   ├── __init__.py
│   │   └── templates.py
│   └── middleware/
│       ├── __init__.py
│       ├── rate_limiter.py
│       └── telemetry.py
└── tests/
    ├── __init__.py
    ├── test_generator.py
    └── conftest.py
```

## Implementation

### 1. Main Application (app/main.py)

```python
"""AI Gateway Service - Main Application"""

from contextlib import asynccontextmanager
from typing import AsyncGenerator

from fastapi import FastAPI, Request, status
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
import structlog

from app.config import settings
from app.middleware.rate_limiter import RateLimiterMiddleware
from app.middleware.telemetry import setup_telemetry
from app.providers import get_provider_registry
from app.services.generator import ServiceGenerator
from app.services.enhancer import WorkflowEnhancer
from app.services.router import InferenceRouter


logger = structlog.get_logger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncGenerator[None, None]:
    """Application lifespan manager."""
    # Startup
    logger.info("Starting AI Gateway service", version=settings.VERSION)
    
    # Initialize providers
    provider_registry = await get_provider_registry()
    app.state.providers = provider_registry
    
    # Initialize services
    app.state.generator = ServiceGenerator(provider_registry)
    app.state.enhancer = WorkflowEnhancer(provider_registry)
    app.state.router = InferenceRouter(provider_registry)
    
    yield
    
    # Shutdown
    logger.info("Shutting down AI Gateway service")
    await provider_registry.close()


app = FastAPI(
    title="OmniRoute AI Gateway",
    description="AI-powered service generation and inference",
    version=settings.VERSION,
    lifespan=lifespan,
)

# Middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.CORS_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
app.add_middleware(RateLimiterMiddleware, redis_url=settings.REDIS_URL)

# Telemetry
setup_telemetry(app)
FastAPIInstrumentor.instrument_app(app)


@app.exception_handler(Exception)
async def global_exception_handler(request: Request, exc: Exception) -> JSONResponse:
    """Global exception handler."""
    logger.error("Unhandled exception", exc_info=exc, path=request.url.path)
    return JSONResponse(
        status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
        content={"error": "internal_error", "message": "An internal error occurred"},
    )


@app.get("/health")
async def health_check():
    """Health check endpoint."""
    return {"status": "healthy", "version": settings.VERSION}


# Import and include routers
from app.api import generation, inference, models

app.include_router(generation.router, prefix="/api/v1", tags=["generation"])
app.include_router(inference.router, prefix="/api/v1", tags=["inference"])
app.include_router(models.router, prefix="/api/v1", tags=["models"])
```

### 2. Service Generator (app/services/generator.py)

```python
"""Service Generation using AI."""

import json
import re
from typing import Any, Optional

import structlog
from pydantic import BaseModel, Field

from app.models.requests import GenerationMode
from app.models.responses import GenerationResult
from app.providers import ProviderRegistry
from app.prompts.templates import (
    SERVICE_GENERATION_SYSTEM,
    WORKFLOW_GENERATION_SYSTEM,
    CODE_GENERATION_SYSTEM,
)


logger = structlog.get_logger(__name__)


class GenerationRequest(BaseModel):
    """Request for service generation."""
    
    tenant_id: str
    mode: GenerationMode
    prompt: str
    context: Optional[dict[str, Any]] = None
    existing_service_id: Optional[str] = None
    use_local_model: bool = False
    temperature: float = Field(default=0.7, ge=0, le=2)
    max_tokens: int = Field(default=4000, ge=100, le=32000)


class ServiceGenerator:
    """AI-powered service generator."""
    
    def __init__(self, providers: ProviderRegistry):
        self.providers = providers
    
    async def generate(self, request: GenerationRequest) -> GenerationResult:
        """Generate a service definition from natural language."""
        log = logger.bind(
            tenant_id=request.tenant_id,
            mode=request.mode.value,
            use_local=request.use_local_model,
        )
        log.info("Starting service generation")
        
        # Select provider and model
        if request.use_local_model:
            provider = self.providers.get("local")
            model = self._select_local_model(request.mode)
        else:
            provider = self.providers.get("anthropic")
            model = "claude-sonnet-4-20250514"
        
        # Build prompts
        system_prompt = self._build_system_prompt(request)
        user_prompt = self._build_user_prompt(request)
        
        # Generate
        try:
            response = await provider.complete(
                model=model,
                system=system_prompt,
                messages=[{"role": "user", "content": user_prompt}],
                temperature=request.temperature,
                max_tokens=request.max_tokens,
            )
        except Exception as e:
            log.error("Generation failed", error=str(e))
            raise
        
        # Parse response
        parsed = self._parse_response(response.content, request.mode)
        
        log.info(
            "Generation completed",
            tokens_used=response.usage.total_tokens,
            model=model,
        )
        
        return GenerationResult(
            status="success",
            service_definition=parsed.get("service_definition"),
            workflow_dsl=parsed.get("workflow_dsl"),
            generated_code=parsed.get("code"),
            explanation=parsed.get("explanation"),
            suggestions=parsed.get("suggestions"),
            tokens_used=response.usage.total_tokens,
            model_used=model,
        )
    
    def _build_system_prompt(self, request: GenerationRequest) -> str:
        """Build the system prompt based on generation mode."""
        base_context = self._get_base_context()
        
        if request.mode == GenerationMode.PROMPT_TO_SERVICE:
            return SERVICE_GENERATION_SYSTEM.format(base_context=base_context)
        elif request.mode == GenerationMode.PROMPT_TO_WORKFLOW:
            return WORKFLOW_GENERATION_SYSTEM.format(base_context=base_context)
        elif request.mode == GenerationMode.CODE_GENERATION:
            return CODE_GENERATION_SYSTEM.format(base_context=base_context)
        else:
            return base_context
    
    def _build_user_prompt(self, request: GenerationRequest) -> str:
        """Build the user prompt with context."""
        parts = [request.prompt]
        
        if request.context:
            parts.append(f"\n\nAdditional context:\n```json\n{json.dumps(request.context, indent=2)}\n```")
        
        if request.existing_service_id:
            parts.append(f"\n\nThis extends existing service: {request.existing_service_id}")
        
        return "\n".join(parts)
    
    def _get_base_context(self) -> str:
        """Get the base context about the platform."""
        return """You are an expert service architect for OmniRoute, a B2B FMCG platform.

You help users create business services that orchestrate:
- Temporal workflows (durable, fault-tolerant execution)
- n8n integrations (400+ pre-built connectors)
- AI capabilities (LLM and local model inference)

Available node types:
1. activity: Execute a Temporal activity (HTTP, database, notification, etc.)
2. ai_action: Call an LLM or local model for AI-powered processing
3. n8n: Execute an n8n workflow for external integrations
4. decision: Conditional branching based on data
5. parallel: Execute multiple branches concurrently
6. wait: Wait for a timer, signal, or condition
7. human_task: Request human approval or input

Always generate valid, well-structured workflow definitions that can be compiled
and executed on the platform."""
    
    def _select_local_model(self, mode: GenerationMode) -> str:
        """Select the best local model for the task."""
        if mode == GenerationMode.CODE_GENERATION:
            return "qwen2.5-coder-32b-instruct"
        return "llama-3.3-70b-instruct"
    
    def _parse_response(
        self, content: str, mode: GenerationMode
    ) -> dict[str, Any]:
        """Parse the AI response to extract structured data."""
        result: dict[str, Any] = {}
        
        # Try to extract JSON from the response
        json_match = re.search(r"```(?:json)?\s*(\{[\s\S]*?\})\s*```", content)
        if not json_match:
            # Try without code blocks
            json_match = re.search(r"(\{[\s\S]*\})", content)
        
        if json_match:
            try:
                parsed_json = json.loads(json_match.group(1))
                
                if mode == GenerationMode.PROMPT_TO_SERVICE:
                    result["service_definition"] = parsed_json
                elif mode == GenerationMode.PROMPT_TO_WORKFLOW:
                    result["workflow_dsl"] = parsed_json
                elif mode == GenerationMode.CODE_GENERATION:
                    result["code"] = parsed_json.get("code", "")
            except json.JSONDecodeError as e:
                logger.warning("Failed to parse JSON from response", error=str(e))
        
        # Extract explanation (text after JSON)
        if json_match:
            explanation_start = json_match.end()
            explanation = content[explanation_start:].strip()
            if explanation:
                result["explanation"] = explanation
                
                # Try to extract suggestions
                suggestions = self._extract_suggestions(explanation)
                if suggestions:
                    result["suggestions"] = suggestions
        else:
            result["explanation"] = content
        
        return result
    
    def _extract_suggestions(self, text: str) -> list[str]:
        """Extract suggestions from the explanation text."""
        suggestions = []
        
        # Look for bullet points or numbered items
        patterns = [
            r"(?:^|\n)\s*[-•]\s*(.+?)(?=\n|$)",
            r"(?:^|\n)\s*\d+\.\s*(.+?)(?=\n|$)",
            r"Suggestion[s]?:\s*(.+?)(?=\n\n|$)",
        ]
        
        for pattern in patterns:
            matches = re.findall(pattern, text, re.MULTILINE)
            for match in matches:
                suggestion = match.strip()
                if len(suggestion) > 10 and len(suggestion) < 200:
                    suggestions.append(suggestion)
        
        return suggestions[:5]  # Limit to 5 suggestions
```

### 3. Provider Interface (app/providers/base.py)

```python
"""Base provider interface."""

from abc import ABC, abstractmethod
from dataclasses import dataclass
from typing import Any, Optional


@dataclass
class Message:
    """Chat message."""
    role: str
    content: str


@dataclass
class Usage:
    """Token usage."""
    prompt_tokens: int
    completion_tokens: int
    total_tokens: int


@dataclass
class CompletionResponse:
    """Completion response from provider."""
    content: str
    model: str
    usage: Usage
    finish_reason: Optional[str] = None


class AIProvider(ABC):
    """Abstract base class for AI providers."""
    
    @abstractmethod
    async def complete(
        self,
        model: str,
        messages: list[dict[str, str]],
        *,
        system: Optional[str] = None,
        temperature: float = 0.7,
        max_tokens: int = 4000,
        **kwargs: Any,
    ) -> CompletionResponse:
        """Generate a completion."""
        pass
    
    @abstractmethod
    async def health_check(self) -> bool:
        """Check if the provider is healthy."""
        pass
    
    @abstractmethod
    async def close(self) -> None:
        """Close any open connections."""
        pass


class ProviderRegistry:
    """Registry of AI providers."""
    
    def __init__(self):
        self._providers: dict[str, AIProvider] = {}
    
    def register(self, name: str, provider: AIProvider) -> None:
        """Register a provider."""
        self._providers[name] = provider
    
    def get(self, name: str) -> AIProvider:
        """Get a provider by name."""
        if name not in self._providers:
            raise ValueError(f"Unknown provider: {name}")
        return self._providers[name]
    
    def list_providers(self) -> list[str]:
        """List registered providers."""
        return list(self._providers.keys())
    
    async def close(self) -> None:
        """Close all providers."""
        for provider in self._providers.values():
            await provider.close()
```

### 4. Anthropic Provider (app/providers/anthropic.py)

```python
"""Anthropic (Claude) provider."""

from typing import Any, Optional

import httpx
import structlog

from app.config import settings
from app.providers.base import AIProvider, CompletionResponse, Usage


logger = structlog.get_logger(__name__)


class AnthropicProvider(AIProvider):
    """Anthropic API provider."""
    
    def __init__(self, api_key: str):
        self.api_key = api_key
        self.base_url = "https://api.anthropic.com/v1"
        self.client = httpx.AsyncClient(
            timeout=120.0,
            headers={
                "x-api-key": api_key,
                "anthropic-version": "2023-06-01",
                "content-type": "application/json",
            },
        )
    
    async def complete(
        self,
        model: str,
        messages: list[dict[str, str]],
        *,
        system: Optional[str] = None,
        temperature: float = 0.7,
        max_tokens: int = 4000,
        **kwargs: Any,
    ) -> CompletionResponse:
        """Generate a completion using Claude."""
        log = logger.bind(model=model, temperature=temperature)
        log.debug("Calling Anthropic API")
        
        payload = {
            "model": model,
            "messages": messages,
            "max_tokens": max_tokens,
            "temperature": temperature,
        }
        
        if system:
            payload["system"] = system
        
        response = await self.client.post(
            f"{self.base_url}/messages",
            json=payload,
        )
        response.raise_for_status()
        
        data = response.json()
        
        content = ""
        for block in data.get("content", []):
            if block.get("type") == "text":
                content += block.get("text", "")
        
        usage = Usage(
            prompt_tokens=data["usage"]["input_tokens"],
            completion_tokens=data["usage"]["output_tokens"],
            total_tokens=data["usage"]["input_tokens"] + data["usage"]["output_tokens"],
        )
        
        log.debug("Anthropic API response received", tokens=usage.total_tokens)
        
        return CompletionResponse(
            content=content,
            model=model,
            usage=usage,
            finish_reason=data.get("stop_reason"),
        )
    
    async def health_check(self) -> bool:
        """Check if Anthropic API is accessible."""
        try:
            response = await self.client.get(f"{self.base_url}/models")
            return response.status_code == 200
        except Exception:
            return False
    
    async def close(self) -> None:
        """Close the HTTP client."""
        await self.client.aclose()
```

### 5. vLLM Local Provider (app/providers/vllm.py)

```python
"""vLLM local inference provider."""

from typing import Any, Optional

import httpx
import structlog

from app.config import settings
from app.providers.base import AIProvider, CompletionResponse, Usage


logger = structlog.get_logger(__name__)


class VLLMProvider(AIProvider):
    """vLLM local inference provider."""
    
    def __init__(self, base_url: str):
        self.base_url = base_url
        self.client = httpx.AsyncClient(timeout=180.0)
    
    async def complete(
        self,
        model: str,
        messages: list[dict[str, str]],
        *,
        system: Optional[str] = None,
        temperature: float = 0.7,
        max_tokens: int = 4000,
        **kwargs: Any,
    ) -> CompletionResponse:
        """Generate a completion using local model."""
        log = logger.bind(model=model, temperature=temperature)
        log.debug("Calling vLLM API")
        
        # Prepend system message if provided
        if system:
            messages = [{"role": "system", "content": system}] + messages
        
        payload = {
            "model": model,
            "messages": messages,
            "max_tokens": max_tokens,
            "temperature": temperature,
            "top_p": 0.95,
        }
        
        response = await self.client.post(
            f"{self.base_url}/v1/chat/completions",
            json=payload,
        )
        response.raise_for_status()
        
        data = response.json()
        
        content = data["choices"][0]["message"]["content"]
        
        usage = Usage(
            prompt_tokens=data["usage"]["prompt_tokens"],
            completion_tokens=data["usage"]["completion_tokens"],
            total_tokens=data["usage"]["total_tokens"],
        )
        
        log.debug("vLLM API response received", tokens=usage.total_tokens)
        
        return CompletionResponse(
            content=content,
            model=model,
            usage=usage,
            finish_reason=data["choices"][0].get("finish_reason"),
        )
    
    async def list_models(self) -> list[str]:
        """List available models on vLLM server."""
        response = await self.client.get(f"{self.base_url}/v1/models")
        response.raise_for_status()
        
        data = response.json()
        return [m["id"] for m in data.get("data", [])]
    
    async def health_check(self) -> bool:
        """Check if vLLM server is healthy."""
        try:
            response = await self.client.get(f"{self.base_url}/health")
            return response.status_code == 200
        except Exception:
            return False
    
    async def close(self) -> None:
        """Close the HTTP client."""
        await self.client.aclose()
```

### 6. API Routes (app/api/generation.py)

```python
"""Generation API routes."""

from fastapi import APIRouter, Request, HTTPException, status
from pydantic import BaseModel, Field
from typing import Any, Optional

from app.models.requests import GenerationMode
from app.services.generator import GenerationRequest


router = APIRouter()


class CreateServiceRequest(BaseModel):
    """Request to generate a service."""
    
    tenant_id: str
    prompt: str
    context: Optional[dict[str, Any]] = None
    use_local_model: bool = False
    temperature: float = Field(default=0.7, ge=0, le=2)
    max_tokens: int = Field(default=4000, ge=100, le=32000)


class EnhanceWorkflowRequest(BaseModel):
    """Request to enhance a workflow."""
    
    tenant_id: str
    workflow_dsl: dict[str, Any]
    enhancement_prompt: str
    temperature: float = Field(default=0.3, ge=0, le=2)


@router.post("/generate/service")
async def generate_service(request: Request, body: CreateServiceRequest):
    """Generate a service from natural language description."""
    generator = request.app.state.generator
    
    gen_request = GenerationRequest(
        tenant_id=body.tenant_id,
        mode=GenerationMode.PROMPT_TO_SERVICE,
        prompt=body.prompt,
        context=body.context,
        use_local_model=body.use_local_model,
        temperature=body.temperature,
        max_tokens=body.max_tokens,
    )
    
    try:
        result = await generator.generate(gen_request)
        return result.model_dump()
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=str(e),
        )


@router.post("/generate/workflow")
async def generate_workflow(request: Request, body: CreateServiceRequest):
    """Generate a workflow DSL from description."""
    generator = request.app.state.generator
    
    gen_request = GenerationRequest(
        tenant_id=body.tenant_id,
        mode=GenerationMode.PROMPT_TO_WORKFLOW,
        prompt=body.prompt,
        context=body.context,
        use_local_model=body.use_local_model,
        temperature=body.temperature,
        max_tokens=body.max_tokens,
    )
    
    try:
        result = await generator.generate(gen_request)
        return result.model_dump()
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=str(e),
        )


@router.post("/enhance/workflow")
async def enhance_workflow(request: Request, body: EnhanceWorkflowRequest):
    """Enhance an existing workflow based on natural language."""
    enhancer = request.app.state.enhancer
    
    try:
        result = await enhancer.enhance(
            tenant_id=body.tenant_id,
            workflow_dsl=body.workflow_dsl,
            enhancement_prompt=body.enhancement_prompt,
            temperature=body.temperature,
        )
        return result.model_dump()
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=str(e),
        )
```

## Deliverables

1. Complete FastAPI application
2. Provider implementations (Anthropic, OpenAI, vLLM)
3. Service generator with prompt templates
4. Workflow enhancer
5. Rate limiting middleware
6. OpenTelemetry integration
7. Comprehensive tests
```

---

# Continue to Part 4 for Frontend Implementation, Infrastructure, and DevOps...
