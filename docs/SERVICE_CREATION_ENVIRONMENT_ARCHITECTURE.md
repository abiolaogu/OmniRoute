# OmniRoute Service Creation Environment (SCE)
## Comprehensive Technical Architecture & Implementation Guide

**Version:** 1.0.0  
**Date:** January 2026  
**Classification:** Technical Specification  

---

## Executive Summary

The OmniRoute Service Creation Environment (SCE) is a **self-service platform** that empowers non-technical users to create, extend, and customize business services through a combination of:

1. **Temporal** - Durable execution engine for fault-tolerant workflows
2. **n8n** - Pre-built integrations and visual automation
3. **AI Engine** - LLM-powered service generation and local model inference

This document provides comprehensive implementation prompts for building the SCE on **Google Cloud Platform (GKE)**, adhering to the **Holy Trinity of Modern Software Development**:

- **Extreme Programming (XP)** - TDD, continuous integration, pair programming, simple design
- **Domain-Driven Design (DDD)** - Bounded contexts, aggregates, domain events, ubiquitous language
- **Legacy Modernization** - Strangler fig pattern, anti-corruption layers, event-driven decoupling

---

## Table of Contents

1. [Strategic Architecture Overview](#1-strategic-architecture-overview)
2. [Domain-Driven Design Model](#2-domain-driven-design-model)
3. [Temporal Workflow Engine](#3-temporal-workflow-engine)
4. [n8n Integration Layer](#4-n8n-integration-layer)
5. [AI Service Generation Engine](#5-ai-service-generation-engine)
6. [Visual Service Designer (WYSIWYG)](#6-visual-service-designer-wysiwyg)
7. [GCP Infrastructure & Deployment](#7-gcp-infrastructure--deployment)
8. [XP Practices Implementation](#8-xp-practices-implementation)
9. [Implementation Prompts](#9-implementation-prompts)
10. [Performance Optimization](#10-performance-optimization)

---

## 1. Strategic Architecture Overview

### 1.1 High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        PRESENTATION LAYER                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐              │
│  │  Visual Service │  │   AI Prompt     │  │   API Gateway   │              │
│  │    Designer     │  │   Interface     │  │    (Kong/Envoy) │              │
│  │  (React Flow)   │  │   (Chat UI)     │  │                 │              │
│  └────────┬────────┘  └────────┬────────┘  └────────┬────────┘              │
└───────────┼─────────────────────┼─────────────────────┼─────────────────────┘
            │                     │                     │
            ▼                     ▼                     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                     SERVICE CREATION LAYER (Go/Rust)                         │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐              │
│  │  Service        │  │   Workflow      │  │   AI Service    │              │
│  │  Registry       │  │   Compiler      │  │   Generator     │              │
│  │  (Bounded Ctx)  │  │   (DSL→Code)    │  │   (LLM+Local)   │              │
│  └────────┬────────┘  └────────┬────────┘  └────────┬────────┘              │
└───────────┼─────────────────────┼─────────────────────┼─────────────────────┘
            │                     │                     │
            ▼                     ▼                     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                     ORCHESTRATION LAYER                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                      TEMPORAL CLUSTER                                │    │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐               │    │
│  │  │   Frontend   │  │   History    │  │   Matching   │               │    │
│  │  │   Service    │  │   Service    │  │   Service    │               │    │
│  │  └──────────────┘  └──────────────┘  └──────────────┘               │    │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                      n8n CLUSTER                                     │    │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐               │    │
│  │  │   Main       │  │   Worker     │  │   Webhook    │               │    │
│  │  │   Instance   │  │   Nodes      │  │   Processor  │               │    │
│  │  └──────────────┘  └──────────────┘  └──────────────┘               │    │
│  └──────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
            │                     │                     │
            ▼                     ▼                     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                     AI INFERENCE LAYER                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐              │
│  │  LLM Gateway    │  │   Local Model   │  │   Model         │              │
│  │  (Claude/GPT)   │  │   Inference     │  │   Registry      │              │
│  │                 │  │   (vLLM/Llama)  │  │   (MLflow)      │              │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘              │
└─────────────────────────────────────────────────────────────────────────────┘
            │                     │                     │
            ▼                     ▼                     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                     PLATFORM SERVICES (OmniRoute Core)                       │
├─────────────────────────────────────────────────────────────────────────────┤
│  Order Service │ Payment Service │ Notification Service │ Gig Platform      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Technology Stack Selection Matrix

| Layer | Technology | Rationale |
|-------|------------|-----------|
| **Visual Designer** | React + React Flow + TypeScript | Best-in-class flow visualization, strong typing |
| **API Gateway** | Kong/Envoy | High performance, extensible, GKE native |
| **Service Registry** | Go | Concurrency, compile-time safety, K8s ecosystem |
| **Workflow Compiler** | Rust | Maximum performance for DSL parsing, memory safety |
| **Temporal Workers** | Go | Official SDK, production-proven, low latency |
| **AI Gateway** | Python (FastAPI) | ML ecosystem, async support |
| **Local Inference** | Rust + vLLM | GPU optimization, low latency |
| **Event Bus** | Kafka + Redpanda | High throughput, exactly-once semantics |
| **Database** | PostgreSQL + TimescaleDB | ACID, time-series for metrics |
| **Cache** | Redis Cluster + Dragonfly | High availability, Dragonfly for large workloads |
| **Search** | Meilisearch | Fast, typo-tolerant, self-hosted |

---

## 2. Domain-Driven Design Model

### 2.1 Bounded Contexts

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        SCE BOUNDED CONTEXTS                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌────────────────────────┐      ┌────────────────────────┐                 │
│  │   SERVICE DEFINITION   │      │   WORKFLOW EXECUTION   │                 │
│  │      CONTEXT           │      │      CONTEXT           │                 │
│  │                        │      │                        │                 │
│  │  Aggregates:           │      │  Aggregates:           │                 │
│  │  • ServiceDefinition   │─────▶│  • WorkflowInstance    │                 │
│  │  • ServiceVersion      │      │  • ExecutionState      │                 │
│  │  • ActivityCatalog     │      │  • ActivityResult      │                 │
│  │                        │      │                        │                 │
│  │  Domain Events:        │      │  Domain Events:        │                 │
│  │  • ServicePublished    │      │  • WorkflowStarted     │                 │
│  │  • ServiceDeprecated   │      │  • ActivityCompleted   │                 │
│  │  • VersionReleased     │      │  • WorkflowFailed      │                 │
│  └────────────────────────┘      └────────────────────────┘                 │
│              │                              │                                │
│              │        Anti-Corruption       │                                │
│              │            Layer             │                                │
│              ▼                              ▼                                │
│  ┌────────────────────────┐      ┌────────────────────────┐                 │
│  │   AI GENERATION        │      │   INTEGRATION          │                 │
│  │      CONTEXT           │      │      CONTEXT           │                 │
│  │                        │      │                        │                 │
│  │  Aggregates:           │      │  Aggregates:           │                 │
│  │  • GenerationRequest   │      │  • Connector           │                 │
│  │  • ModelConfiguration  │      │  • CredentialVault     │                 │
│  │  • PromptTemplate      │      │  • WebhookEndpoint     │                 │
│  │                        │      │                        │                 │
│  │  Domain Events:        │      │  Domain Events:        │                 │
│  │  • ServiceGenerated    │      │  • ConnectorConfigured │                 │
│  │  • GenerationFailed    │      │  • WebhookTriggered    │                 │
│  │  • ModelSwitched       │      │  • IntegrationFailed   │                 │
│  └────────────────────────┘      └────────────────────────┘                 │
│                                                                              │
│  ┌────────────────────────┐      ┌────────────────────────┐                 │
│  │   TENANT MANAGEMENT    │      │   OBSERVABILITY        │                 │
│  │      CONTEXT           │      │      CONTEXT           │                 │
│  │                        │      │                        │                 │
│  │  Aggregates:           │      │  Aggregates:           │                 │
│  │  • Tenant              │      │  • MetricStream        │                 │
│  │  • UsageQuota          │      │  • AuditLog            │                 │
│  │  • BillingAccount      │      │  • AlertRule           │                 │
│  └────────────────────────┘      └────────────────────────┘                 │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Aggregate Design - Service Definition

```go
// Domain Model - Service Definition Aggregate Root
package domain

import (
    "time"
    "github.com/google/uuid"
)

// ServiceDefinition is the Aggregate Root for service creation
type ServiceDefinition struct {
    // Identity
    ID        ServiceID   `json:"id"`
    TenantID  TenantID    `json:"tenant_id"`
    
    // Core Properties
    Name        ServiceName        `json:"name"`
    Description string             `json:"description"`
    Category    ServiceCategory    `json:"category"`
    
    // Workflow Definition (Value Object)
    Workflow    WorkflowGraph      `json:"workflow"`
    
    // Versioning
    Versions    []ServiceVersion   `json:"versions"`
    ActiveVersion VersionID        `json:"active_version"`
    
    // Lifecycle
    Status      ServiceStatus      `json:"status"`
    CreatedAt   time.Time          `json:"created_at"`
    UpdatedAt   time.Time          `json:"updated_at"`
    PublishedAt *time.Time         `json:"published_at,omitempty"`
    
    // Domain Events (transient)
    events      []DomainEvent
}

// Value Objects
type ServiceID uuid.UUID
type TenantID uuid.UUID
type VersionID uuid.UUID

type ServiceName string
func (n ServiceName) Validate() error {
    if len(n) < 3 || len(n) > 100 {
        return ErrInvalidServiceName
    }
    return nil
}

type ServiceStatus string
const (
    ServiceStatusDraft     ServiceStatus = "draft"
    ServiceStatusPublished ServiceStatus = "published"
    ServiceStatusDeprecated ServiceStatus = "deprecated"
    ServiceStatusArchived   ServiceStatus = "archived"
)

type ServiceCategory string
const (
    CategoryAutomation    ServiceCategory = "automation"
    CategoryIntegration   ServiceCategory = "integration"
    CategoryNotification  ServiceCategory = "notification"
    CategoryPayment       ServiceCategory = "payment"
    CategoryFulfillment   ServiceCategory = "fulfillment"
    CategoryAnalytics     ServiceCategory = "analytics"
    CategoryCustom        ServiceCategory = "custom"
)

// WorkflowGraph is a Value Object representing the visual workflow
type WorkflowGraph struct {
    Nodes       []WorkflowNode     `json:"nodes"`
    Edges       []WorkflowEdge     `json:"edges"`
    Variables   []WorkflowVariable `json:"variables"`
    Triggers    []WorkflowTrigger  `json:"triggers"`
    ErrorPolicy ErrorHandlingPolicy `json:"error_policy"`
}

type WorkflowNode struct {
    ID          string                 `json:"id"`
    Type        NodeType               `json:"type"`
    ActivityRef *ActivityReference     `json:"activity_ref,omitempty"`
    AIAction    *AIActionConfig        `json:"ai_action,omitempty"`
    n8nNode     *N8NNodeConfig         `json:"n8n_node,omitempty"`
    Position    Position               `json:"position"`
    Config      map[string]interface{} `json:"config"`
}

type NodeType string
const (
    NodeTypeActivity   NodeType = "activity"      // Temporal Activity
    NodeTypeSubflow    NodeType = "subflow"       // Child Workflow
    NodeTypeAIAction   NodeType = "ai_action"     // LLM/Local Model
    NodeTypeN8N        NodeType = "n8n"           // n8n Integration
    NodeTypeDecision   NodeType = "decision"      // Conditional Branch
    NodeTypeParallel   NodeType = "parallel"      // Parallel Execution
    NodeTypeWait       NodeType = "wait"          // Timer/Signal Wait
    NodeTypeHuman      NodeType = "human_task"    // Human-in-the-loop
)

// Domain Methods
func (s *ServiceDefinition) Publish() error {
    if s.Status == ServiceStatusArchived {
        return ErrCannotPublishArchived
    }
    if len(s.Workflow.Nodes) == 0 {
        return ErrEmptyWorkflow
    }
    
    s.Status = ServiceStatusPublished
    now := time.Now().UTC()
    s.PublishedAt = &now
    s.UpdatedAt = now
    
    s.events = append(s.events, ServicePublishedEvent{
        ServiceID:   s.ID,
        TenantID:    s.TenantID,
        Version:     s.ActiveVersion,
        PublishedAt: now,
    })
    
    return nil
}

func (s *ServiceDefinition) AddVersion(graph WorkflowGraph, notes string) (VersionID, error) {
    version := ServiceVersion{
        ID:           VersionID(uuid.New()),
        ServiceID:    s.ID,
        VersionNumber: s.nextVersionNumber(),
        Workflow:     graph,
        ReleaseNotes: notes,
        CreatedAt:    time.Now().UTC(),
    }
    
    s.Versions = append(s.Versions, version)
    s.ActiveVersion = version.ID
    
    s.events = append(s.events, VersionReleasedEvent{
        ServiceID: s.ID,
        VersionID: version.ID,
        Number:    version.VersionNumber,
    })
    
    return version.ID, nil
}

func (s *ServiceDefinition) PullEvents() []DomainEvent {
    events := s.events
    s.events = nil
    return events
}
```

### 2.3 Context Mapping

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        CONTEXT MAP                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│     Service Definition ◄────────────────────► Workflow Execution             │
│           Context              Published/       Context                      │
│                               Subscriber                                     │
│               │                                       │                      │
│               │                                       │                      │
│     Anti-Corruption                          Anti-Corruption                 │
│         Layer                                    Layer                       │
│               │                                       │                      │
│               ▼                                       ▼                      │
│      AI Generation ◄─────────────────────────► Integration                   │
│           Context        Shared Kernel          Context                      │
│                         (Event Schemas)                                      │
│                                                                              │
│                                                                              │
│     ┌─────────────────────────────────────────────────────────────┐         │
│     │              UPSTREAM SYSTEMS (OmniRoute Core)               │         │
│     │                                                              │         │
│     │  Order Service ────► Conformist                              │         │
│     │  Payment Service ──► Conformist                              │         │
│     │  Notification ─────► Conformist                              │         │
│     │  Gig Platform ─────► Conformist                              │         │
│     └─────────────────────────────────────────────────────────────┘         │
│                                                                              │
│                                                                              │
│     ┌─────────────────────────────────────────────────────────────┐         │
│     │              EXTERNAL SYSTEMS (ACL Required)                 │         │
│     │                                                              │         │
│     │  n8n ─────────────► Anti-Corruption Layer                    │         │
│     │  Temporal ─────────► Anti-Corruption Layer                   │         │
│     │  LLM APIs ─────────► Anti-Corruption Layer                   │         │
│     │  Local Models ─────► Anti-Corruption Layer                   │         │
│     └─────────────────────────────────────────────────────────────┘         │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Temporal Workflow Engine

### 3.1 Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     TEMPORAL INTEGRATION ARCHITECTURE                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    TEMPORAL SERVER (GKE)                             │    │
│  │                                                                      │    │
│  │  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐                 │    │
│  │  │  Frontend    │ │   History    │ │   Matching   │                 │    │
│  │  │  (gRPC API)  │ │   Service    │ │   Service    │                 │    │
│  │  └──────────────┘ └──────────────┘ └──────────────┘                 │    │
│  │         │                │                │                          │    │
│  │         └────────────────┼────────────────┘                          │    │
│  │                          │                                           │    │
│  │                   ┌──────▼──────┐                                    │    │
│  │                   │   Cassandra │   (or PostgreSQL)                  │    │
│  │                   │   Cluster   │                                    │    │
│  │                   └─────────────┘                                    │    │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                          │                                                   │
│                          │  gRPC                                            │
│                          ▼                                                   │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                    WORKER FLEET (Auto-scaled)                         │   │
│  │                                                                       │   │
│  │  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐       │   │
│  │  │  Core Workflow  │  │  Integration    │  │  AI Inference   │       │   │
│  │  │  Workers (Go)   │  │  Workers (Go)   │  │  Workers (Py)   │       │   │
│  │  │                 │  │                 │  │                 │       │   │
│  │  │  • UserService  │  │  • n8nActivity  │  │  • LLMActivity  │       │   │
│  │  │  • OrderFlow    │  │  • WebhookAct   │  │  • LocalModel   │       │   │
│  │  │  • PaymentFlow  │  │  • APICallAct   │  │  • Embedding    │       │   │
│  │  └─────────────────┘  └─────────────────┘  └─────────────────┘       │   │
│  │                                                                       │   │
│  │  Task Queues:                                                         │   │
│  │  • omniroute-core        (high priority)                              │   │
│  │  • omniroute-integration (medium priority)                            │   │
│  │  • omniroute-ai          (low priority, GPU workers)                  │   │
│  │  • omniroute-batch       (background processing)                      │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 DSL to Temporal Code Generation

The **Workflow Compiler** transforms the visual DSL (JSON) into executable Temporal workflow code.

```rust
// workflow-compiler/src/lib.rs
// Rust-based DSL compiler for maximum performance

use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Debug, Deserialize)]
pub struct WorkflowDSL {
    pub id: String,
    pub name: String,
    pub tenant_id: String,
    pub nodes: Vec<Node>,
    pub edges: Vec<Edge>,
    pub variables: Vec<Variable>,
    pub triggers: Vec<Trigger>,
}

#[derive(Debug, Deserialize)]
pub struct Node {
    pub id: String,
    pub node_type: NodeType,
    pub config: HashMap<String, serde_json::Value>,
    pub retry_policy: Option<RetryPolicy>,
    pub timeout: Option<Duration>,
}

#[derive(Debug, Deserialize)]
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

pub struct TemporalCodeGenerator {
    templates: HashMap<NodeType, String>,
}

impl TemporalCodeGenerator {
    pub fn generate(&self, dsl: &WorkflowDSL) -> Result<GeneratedWorkflow, CompilerError> {
        let mut code = String::new();
        
        // Generate workflow function signature
        code.push_str(&self.generate_workflow_header(dsl));
        
        // Build execution graph
        let graph = self.build_execution_graph(&dsl.nodes, &dsl.edges)?;
        
        // Generate code for each node in topological order
        for node in graph.topological_sort() {
            code.push_str(&self.generate_node_code(node, dsl)?);
        }
        
        // Generate workflow footer
        code.push_str(&self.generate_workflow_footer(dsl));
        
        Ok(GeneratedWorkflow {
            workflow_code: code,
            activity_stubs: self.extract_activity_stubs(dsl),
            required_workers: self.determine_required_workers(dsl),
        })
    }
    
    fn generate_node_code(&self, node: &Node, dsl: &WorkflowDSL) -> Result<String, CompilerError> {
        match node.node_type {
            NodeType::Activity => self.gen_activity_call(node),
            NodeType::AiAction => self.gen_ai_action(node),
            NodeType::N8n => self.gen_n8n_call(node),
            NodeType::Decision => self.gen_decision_branch(node, dsl),
            NodeType::Parallel => self.gen_parallel_execution(node, dsl),
            NodeType::Wait => self.gen_wait_statement(node),
            NodeType::HumanTask => self.gen_human_task(node),
            NodeType::Subflow => self.gen_child_workflow(node),
        }
    }
    
    fn gen_activity_call(&self, node: &Node) -> Result<String, CompilerError> {
        let activity_name = node.config.get("activity_name")
            .ok_or(CompilerError::MissingConfig("activity_name"))?;
        let input_mapping = node.config.get("input_mapping")
            .map(|v| self.generate_input_mapping(v))
            .unwrap_or_default();
        
        Ok(format!(r#"
    // Activity: {node_id}
    var {result_var} {result_type}
    err = workflow.ExecuteActivity(ctx, {activity_name}, {input}).Get(ctx, &{result_var})
    if err != nil {{
        return nil, err
    }}
"#,
            node_id = node.id,
            result_var = self.result_var_name(&node.id),
            result_type = self.infer_result_type(node),
            activity_name = activity_name,
            input = input_mapping,
        ))
    }
    
    fn gen_ai_action(&self, node: &Node) -> Result<String, CompilerError> {
        let model = node.config.get("model").unwrap_or(&serde_json::json!("claude-3-sonnet"));
        let prompt_template = node.config.get("prompt_template")
            .ok_or(CompilerError::MissingConfig("prompt_template"))?;
        let use_local = node.config.get("use_local_model").and_then(|v| v.as_bool()).unwrap_or(false);
        
        if use_local {
            Ok(format!(r#"
    // AI Action (Local): {node_id}
    aiInput := &ai.LocalInferenceInput{{
        Model: "{model}",
        Prompt: fmt.Sprintf(`{prompt}`, {vars}),
        MaxTokens: {max_tokens},
    }}
    var aiResult ai.InferenceResult
    err = workflow.ExecuteActivity(
        workflow.WithTaskQueue(ctx, "omniroute-ai"),
        activities.LocalModelInference,
        aiInput,
    ).Get(ctx, &aiResult)
    if err != nil {{
        return nil, err
    }}
"#,
                node_id = node.id,
                model = model,
                prompt = prompt_template,
                vars = self.extract_template_vars(prompt_template),
                max_tokens = node.config.get("max_tokens").unwrap_or(&serde_json::json!(1000)),
            ))
        } else {
            Ok(format!(r#"
    // AI Action (LLM): {node_id}
    aiInput := &ai.LLMInput{{
        Provider: "{provider}",
        Model: "{model}",
        Messages: []ai.Message{{
            {{Role: "user", Content: fmt.Sprintf(`{prompt}`, {vars})}},
        }},
    }}
    var aiResult ai.LLMResult
    err = workflow.ExecuteActivity(ctx, activities.CallLLM, aiInput).Get(ctx, &aiResult)
    if err != nil {{
        return nil, err
    }}
"#,
                node_id = node.id,
                provider = node.config.get("provider").unwrap_or(&serde_json::json!("anthropic")),
                model = model,
                prompt = prompt_template,
                vars = self.extract_template_vars(prompt_template),
            ))
        }
    }
    
    fn gen_n8n_call(&self, node: &Node) -> Result<String, CompilerError> {
        let workflow_id = node.config.get("n8n_workflow_id")
            .ok_or(CompilerError::MissingConfig("n8n_workflow_id"))?;
        let input_data = node.config.get("input_data")
            .map(|v| serde_json::to_string(v).unwrap())
            .unwrap_or_else(|| "nil".to_string());
        
        Ok(format!(r#"
    // n8n Integration: {node_id}
    n8nInput := &integration.N8NExecutionInput{{
        WorkflowID: "{workflow_id}",
        InputData: {input_data},
        WaitForCompletion: true,
    }}
    var n8nResult integration.N8NExecutionResult
    err = workflow.ExecuteActivity(
        workflow.WithTaskQueue(ctx, "omniroute-integration"),
        activities.ExecuteN8NWorkflow,
        n8nInput,
    ).Get(ctx, &n8nResult)
    if err != nil {{
        return nil, err
    }}
"#,
            node_id = node.id,
            workflow_id = workflow_id,
            input_data = input_data,
        ))
    }
}
```

### 3.3 Temporal Worker Implementation

```go
// workers/core/main.go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "go.temporal.io/sdk/client"
    "go.temporal.io/sdk/worker"
    "go.opentelemetry.io/otel"
    
    "github.com/omniroute/sce/internal/activities"
    "github.com/omniroute/sce/internal/workflows"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // Initialize OpenTelemetry
    tp, err := initTracer()
    if err != nil {
        log.Fatal(err)
    }
    defer tp.Shutdown(ctx)
    otel.SetTracerProvider(tp)
    
    // Create Temporal client with interceptors
    c, err := client.Dial(client.Options{
        HostPort:  os.Getenv("TEMPORAL_HOST"),
        Namespace: os.Getenv("TEMPORAL_NAMESPACE"),
        Interceptors: []client.Interceptor{
            temporal.NewTracingInterceptor(temporal.TracerOptions{}),
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()
    
    // Initialize activity dependencies
    deps := activities.NewDependencies(
        activities.WithPostgres(mustInitPostgres()),
        activities.WithRedis(mustInitRedis()),
        activities.WithKafka(mustInitKafka()),
        activities.WithN8NClient(mustInitN8NClient()),
        activities.WithAIGateway(mustInitAIGateway()),
    )
    
    // Create worker
    w := worker.New(c, "omniroute-core", worker.Options{
        MaxConcurrentActivityExecutionSize:     100,
        MaxConcurrentWorkflowTaskExecutionSize: 50,
        WorkerStopTimeout:                      30 * time.Second,
    })
    
    // Register workflows
    w.RegisterWorkflow(workflows.UserDefinedServiceWorkflow)
    w.RegisterWorkflow(workflows.OrderFulfillmentWorkflow)
    w.RegisterWorkflow(workflows.PaymentProcessingWorkflow)
    w.RegisterWorkflow(workflows.NotificationOrchestrationWorkflow)
    
    // Register activities
    w.RegisterActivity(activities.NewCoreActivities(deps))
    w.RegisterActivity(activities.NewIntegrationActivities(deps))
    
    // Start worker
    go func() {
        if err := w.Run(worker.InterruptCh()); err != nil {
            log.Fatal(err)
        }
    }()
    
    // Graceful shutdown
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    <-sigCh
    
    log.Println("Shutting down worker...")
    w.Stop()
}
```

---

## 4. n8n Integration Layer

### 4.1 n8n Cluster Architecture

```yaml
# k8s/n8n/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: n8n-main
  namespace: omniroute-sce
spec:
  replicas: 2
  selector:
    matchLabels:
      app: n8n
      component: main
  template:
    metadata:
      labels:
        app: n8n
        component: main
    spec:
      containers:
      - name: n8n
        image: n8nio/n8n:latest
        ports:
        - containerPort: 5678
        env:
        - name: N8N_ENCRYPTION_KEY
          valueFrom:
            secretKeyRef:
              name: n8n-secrets
              key: encryption-key
        - name: DB_TYPE
          value: "postgresdb"
        - name: DB_POSTGRESDB_HOST
          value: "postgres-sce.omniroute-sce.svc.cluster.local"
        - name: DB_POSTGRESDB_DATABASE
          value: "n8n"
        - name: EXECUTIONS_MODE
          value: "queue"  # Enable queue mode for scalability
        - name: QUEUE_BULL_REDIS_HOST
          value: "redis-sce.omniroute-sce.svc.cluster.local"
        - name: N8N_DIAGNOSTICS_ENABLED
          value: "false"
        - name: N8N_HIRING_BANNER_ENABLED
          value: "false"
        - name: WEBHOOK_URL
          value: "https://n8n.omniroute.io"
        resources:
          requests:
            memory: "512Mi"
            cpu: "250m"
          limits:
            memory: "2Gi"
            cpu: "1000m"
        livenessProbe:
          httpGet:
            path: /healthz
            port: 5678
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /healthz
            port: 5678
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: n8n-worker
  namespace: omniroute-sce
spec:
  replicas: 5
  selector:
    matchLabels:
      app: n8n
      component: worker
  template:
    metadata:
      labels:
        app: n8n
        component: worker
    spec:
      containers:
      - name: n8n-worker
        image: n8nio/n8n:latest
        command: ["n8n", "worker"]
        env:
        # ... same env vars as main
        - name: EXECUTIONS_MODE
          value: "queue"
        resources:
          requests:
            memory: "256Mi"
            cpu: "100m"
          limits:
            memory: "1Gi"
            cpu: "500m"
```

### 4.2 n8n Integration Gateway

```go
// internal/integration/n8n/gateway.go
package n8n

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    
    "github.com/sony/gobreaker"
    "go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("sce/n8n-gateway")

// Gateway provides integration with n8n instance
type Gateway struct {
    baseURL    string
    apiKey     string
    httpClient *http.Client
    breaker    *gobreaker.CircuitBreaker
}

type GatewayConfig struct {
    BaseURL          string
    APIKey           string
    Timeout          time.Duration
    MaxIdleConns     int
    BreakerThreshold uint32
}

func NewGateway(cfg GatewayConfig) *Gateway {
    transport := &http.Transport{
        MaxIdleConns:        cfg.MaxIdleConns,
        MaxIdleConnsPerHost: cfg.MaxIdleConns,
        IdleConnTimeout:     90 * time.Second,
    }
    
    breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
        Name:        "n8n-gateway",
        MaxRequests: 5,
        Interval:    10 * time.Second,
        Timeout:     30 * time.Second,
        ReadyToTrip: func(counts gobreaker.Counts) bool {
            return counts.ConsecutiveFailures > cfg.BreakerThreshold
        },
    })
    
    return &Gateway{
        baseURL: cfg.BaseURL,
        apiKey:  cfg.APIKey,
        httpClient: &http.Client{
            Transport: transport,
            Timeout:   cfg.Timeout,
        },
        breaker: breaker,
    }
}

// ExecuteWorkflow triggers an n8n workflow and optionally waits for completion
type ExecuteWorkflowInput struct {
    WorkflowID        string                 `json:"workflow_id"`
    InputData         map[string]interface{} `json:"input_data"`
    WaitForCompletion bool                   `json:"wait_for_completion"`
    WebhookPath       string                 `json:"webhook_path,omitempty"`
}

type ExecuteWorkflowResult struct {
    ExecutionID string                 `json:"execution_id"`
    Status      string                 `json:"status"`
    Data        map[string]interface{} `json:"data,omitempty"`
    StartedAt   time.Time              `json:"started_at"`
    FinishedAt  *time.Time             `json:"finished_at,omitempty"`
}

func (g *Gateway) ExecuteWorkflow(ctx context.Context, input ExecuteWorkflowInput) (*ExecuteWorkflowResult, error) {
    ctx, span := tracer.Start(ctx, "n8n.ExecuteWorkflow")
    defer span.End()
    
    // Use circuit breaker
    result, err := g.breaker.Execute(func() (interface{}, error) {
        if input.WebhookPath != "" {
            return g.triggerWebhook(ctx, input)
        }
        return g.triggerWorkflowAPI(ctx, input)
    })
    
    if err != nil {
        span.RecordError(err)
        return nil, fmt.Errorf("execute n8n workflow: %w", err)
    }
    
    return result.(*ExecuteWorkflowResult), nil
}

func (g *Gateway) triggerWebhook(ctx context.Context, input ExecuteWorkflowInput) (*ExecuteWorkflowResult, error) {
    url := fmt.Sprintf("%s/webhook/%s", g.baseURL, input.WebhookPath)
    
    body, err := json.Marshal(input.InputData)
    if err != nil {
        return nil, err
    }
    
    req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := g.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode >= 400 {
        return nil, fmt.Errorf("n8n webhook returned status %d", resp.StatusCode)
    }
    
    var result ExecuteWorkflowResult
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    return &result, nil
}

func (g *Gateway) triggerWorkflowAPI(ctx context.Context, input ExecuteWorkflowInput) (*ExecuteWorkflowResult, error) {
    url := fmt.Sprintf("%s/api/v1/workflows/%s/execute", g.baseURL, input.WorkflowID)
    
    payload := map[string]interface{}{
        "data": input.InputData,
    }
    body, err := json.Marshal(payload)
    if err != nil {
        return nil, err
    }
    
    req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-N8N-API-KEY", g.apiKey)
    
    resp, err := g.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode >= 400 {
        return nil, fmt.Errorf("n8n API returned status %d", resp.StatusCode)
    }
    
    var apiResp struct {
        Data struct {
            ExecutionID string `json:"executionId"`
        } `json:"data"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
        return nil, err
    }
    
    result := &ExecuteWorkflowResult{
        ExecutionID: apiResp.Data.ExecutionID,
        Status:      "running",
        StartedAt:   time.Now(),
    }
    
    // Wait for completion if requested
    if input.WaitForCompletion {
        return g.waitForExecution(ctx, result.ExecutionID)
    }
    
    return result, nil
}

func (g *Gateway) waitForExecution(ctx context.Context, executionID string) (*ExecuteWorkflowResult, error) {
    ticker := time.NewTicker(500 * time.Millisecond)
    defer ticker.Stop()
    
    timeout := time.After(5 * time.Minute)
    
    for {
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        case <-timeout:
            return nil, fmt.Errorf("execution timeout")
        case <-ticker.C:
            result, err := g.getExecutionStatus(ctx, executionID)
            if err != nil {
                return nil, err
            }
            if result.Status == "success" || result.Status == "error" {
                return result, nil
            }
        }
    }
}

// ListWorkflows returns available n8n workflows for the tenant
func (g *Gateway) ListWorkflows(ctx context.Context, tenantID string) ([]WorkflowMeta, error) {
    ctx, span := tracer.Start(ctx, "n8n.ListWorkflows")
    defer span.End()
    
    url := fmt.Sprintf("%s/api/v1/workflows?tags=%s", g.baseURL, tenantID)
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("X-N8N-API-KEY", g.apiKey)
    
    resp, err := g.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var apiResp struct {
        Data []struct {
            ID        string `json:"id"`
            Name      string `json:"name"`
            Active    bool   `json:"active"`
            CreatedAt string `json:"createdAt"`
            UpdatedAt string `json:"updatedAt"`
        } `json:"data"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
        return nil, err
    }
    
    workflows := make([]WorkflowMeta, len(apiResp.Data))
    for i, w := range apiResp.Data {
        workflows[i] = WorkflowMeta{
            ID:        w.ID,
            Name:      w.Name,
            Active:    w.Active,
            CreatedAt: w.CreatedAt,
            UpdatedAt: w.UpdatedAt,
        }
    }
    
    return workflows, nil
}

type WorkflowMeta struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Active      bool   `json:"active"`
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`
}
```

---

## 5. AI Service Generation Engine

### 5.1 Multi-Model Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     AI SERVICE GENERATION ENGINE                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    AI GATEWAY (FastAPI + Python)                     │    │
│  │                                                                      │    │
│  │  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐                 │    │
│  │  │   Router     │ │   Rate       │ │   Model      │                 │    │
│  │  │   (LLM/Local)│ │   Limiter    │ │   Registry   │                 │    │
│  │  └──────────────┘ └──────────────┘ └──────────────┘                 │    │
│  │         │                │                │                          │    │
│  │         └────────────────┼────────────────┘                          │    │
│  │                          │                                           │    │
│  │  ┌───────────────────────▼───────────────────────┐                   │    │
│  │  │              INFERENCE ROUTER                  │                   │    │
│  │  │                                                │                   │    │
│  │  │  ┌─────────┐ ┌─────────┐ ┌─────────┐         │                   │    │
│  │  │  │ Claude  │ │ GPT-4o  │ │ Gemini  │  Cloud  │                   │    │
│  │  │  └─────────┘ └─────────┘ └─────────┘         │                   │    │
│  │  │                                                │                   │    │
│  │  │  ┌─────────┐ ┌─────────┐ ┌─────────┐         │                   │    │
│  │  │  │ Llama   │ │ Mistral │ │ Qwen    │  Local  │                   │    │
│  │  │  └─────────┘ └─────────┘ └─────────┘         │                   │    │
│  │  └────────────────────────────────────────────────┘                  │    │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                          │                                                   │
│                          ▼                                                   │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                    LOCAL INFERENCE CLUSTER                            │   │
│  │                                                                       │   │
│  │  ┌─────────────────────────────────────────────────────────────┐     │   │
│  │  │                      vLLM Inference Server                   │     │   │
│  │  │                                                              │     │   │
│  │  │  GPU Pool: 4x NVIDIA A100 (GKE Autopilot GPU Nodes)         │     │   │
│  │  │                                                              │     │   │
│  │  │  Models Loaded:                                              │     │   │
│  │  │  • llama-3.3-70b-instruct (primary)                         │     │   │
│  │  │  • mistral-large-instruct-2411 (fallback)                   │     │   │
│  │  │  • qwen2.5-coder-32b (code generation)                      │     │   │
│  │  │                                                              │     │   │
│  │  │  Features:                                                   │     │   │
│  │  │  • PagedAttention for efficient KV cache                    │     │   │
│  │  │  • Continuous batching                                       │     │   │
│  │  │  • Tensor parallelism across GPUs                           │     │   │
│  │  └─────────────────────────────────────────────────────────────┘     │   │
│  │                                                                       │   │
│  │  ┌─────────────────────────────────────────────────────────────┐     │   │
│  │  │                      Embedding Service                       │     │   │
│  │  │                                                              │     │   │
│  │  │  • sentence-transformers/all-MiniLM-L6-v2                   │     │   │
│  │  │  • BGE-M3 (multilingual)                                    │     │   │
│  │  └─────────────────────────────────────────────────────────────┘     │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 5.2 AI Gateway Implementation

```python
# ai-gateway/app/main.py
from fastapi import FastAPI, HTTPException, Depends, BackgroundTasks
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel, Field
from typing import Optional, List, Dict, Any, Literal
from enum import Enum
import httpx
import asyncio
from opentelemetry import trace
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor

from app.config import settings
from app.models import (
    ServiceGenerationRequest,
    ServiceGenerationResponse,
    InferenceRequest,
    InferenceResponse,
)
from app.providers import claude, openai, local_vllm
from app.rate_limiter import RateLimiter
from app.model_registry import ModelRegistry

app = FastAPI(title="OmniRoute AI Gateway", version="1.0.0")
FastAPIInstrumentor.instrument_app(app)
tracer = trace.get_tracer(__name__)

# Initialize components
rate_limiter = RateLimiter(redis_url=settings.REDIS_URL)
model_registry = ModelRegistry()


class ModelProvider(str, Enum):
    ANTHROPIC = "anthropic"
    OPENAI = "openai"
    LOCAL = "local"


class GenerationMode(str, Enum):
    PROMPT_TO_SERVICE = "prompt_to_service"
    PROMPT_TO_WORKFLOW = "prompt_to_workflow"
    CODE_GENERATION = "code_generation"
    DOCUMENTATION = "documentation"


# Service Generation Models
class ServiceGenerationRequest(BaseModel):
    tenant_id: str
    mode: GenerationMode
    prompt: str
    context: Optional[Dict[str, Any]] = None
    existing_service_id: Optional[str] = None
    preferred_provider: Optional[ModelProvider] = None
    use_local_model: bool = False
    temperature: float = Field(default=0.7, ge=0, le=2)
    max_tokens: int = Field(default=4000, ge=100, le=32000)


class ServiceGenerationResponse(BaseModel):
    request_id: str
    status: Literal["success", "error", "pending"]
    service_definition: Optional[Dict[str, Any]] = None
    workflow_dsl: Optional[Dict[str, Any]] = None
    generated_code: Optional[str] = None
    explanation: Optional[str] = None
    suggestions: Optional[List[str]] = None
    tokens_used: int
    model_used: str
    latency_ms: int


@app.post("/api/v1/generate/service", response_model=ServiceGenerationResponse)
async def generate_service(
    request: ServiceGenerationRequest,
    background_tasks: BackgroundTasks,
):
    """
    Generate a service definition from natural language prompt.
    
    This endpoint uses either cloud LLMs (Claude/GPT-4) or local models (Llama/Mistral)
    to convert a user's natural language description into a complete service definition
    that can be deployed on the OmniRoute platform.
    """
    with tracer.start_as_current_span("generate_service") as span:
        span.set_attribute("tenant_id", request.tenant_id)
        span.set_attribute("mode", request.mode.value)
        
        # Rate limiting check
        if not await rate_limiter.check(request.tenant_id, "service_generation"):
            raise HTTPException(status_code=429, detail="Rate limit exceeded")
        
        # Build the generation prompt
        system_prompt = build_system_prompt(request.mode)
        user_prompt = build_user_prompt(request)
        
        # Select model based on request
        if request.use_local_model:
            provider = local_vllm
            model = select_local_model(request.mode)
        elif request.preferred_provider == ModelProvider.OPENAI:
            provider = openai
            model = "gpt-4o"
        else:
            provider = claude
            model = "claude-sonnet-4-20250514"
        
        span.set_attribute("model", model)
        
        # Execute inference
        start_time = asyncio.get_event_loop().time()
        
        try:
            response = await provider.complete(
                model=model,
                system=system_prompt,
                messages=[{"role": "user", "content": user_prompt}],
                temperature=request.temperature,
                max_tokens=request.max_tokens,
            )
        except Exception as e:
            span.record_exception(e)
            raise HTTPException(status_code=500, detail=str(e))
        
        latency_ms = int((asyncio.get_event_loop().time() - start_time) * 1000)
        
        # Parse the generated content
        parsed = parse_generation_response(response.content, request.mode)
        
        # Log usage for billing
        background_tasks.add_task(
            log_usage,
            tenant_id=request.tenant_id,
            model=model,
            tokens=response.usage.total_tokens,
        )
        
        return ServiceGenerationResponse(
            request_id=generate_request_id(),
            status="success",
            service_definition=parsed.get("service_definition"),
            workflow_dsl=parsed.get("workflow_dsl"),
            generated_code=parsed.get("code"),
            explanation=parsed.get("explanation"),
            suggestions=parsed.get("suggestions"),
            tokens_used=response.usage.total_tokens,
            model_used=model,
            latency_ms=latency_ms,
        )


def build_system_prompt(mode: GenerationMode) -> str:
    """Build specialized system prompts for different generation modes."""
    
    base_context = """You are an expert service architect for the OmniRoute B2B FMCG platform.
You help users create business services, workflows, and integrations using a combination of:
- Temporal (durable workflow execution)
- n8n (pre-built integrations with 400+ services)
- Custom activities (HTTP calls, database operations, notifications)

The platform uses Domain-Driven Design principles with the following bounded contexts:
- Service Definition Context
- Workflow Execution Context
- AI Generation Context
- Integration Context
- Tenant Management Context
"""

    if mode == GenerationMode.PROMPT_TO_SERVICE:
        return base_context + """
When generating a service definition, output a JSON object with this structure:
{
  "name": "ServiceName",
  "description": "...",
  "category": "automation|integration|notification|payment|fulfillment|analytics|custom",
  "workflow": {
    "nodes": [...],
    "edges": [...],
    "triggers": [...],
    "variables": [...]
  },
  "inputs": [...],
  "outputs": [...],
  "error_handling": {...}
}

Each node should have:
- id: unique identifier
- type: activity|subflow|ai_action|n8n|decision|parallel|wait|human_task
- config: type-specific configuration

Followed by a plain-text explanation of what the service does and any suggestions.
"""
    
    elif mode == GenerationMode.PROMPT_TO_WORKFLOW:
        return base_context + """
Generate a Temporal-compatible workflow definition in JSON DSL format.
Focus on the workflow structure, retry policies, and activity configurations.
Include both the DSL and generated Go code for the workflow.
"""
    
    elif mode == GenerationMode.CODE_GENERATION:
        return base_context + """
Generate production-ready Go code for Temporal workflows and activities.
Follow these principles:
- Use proper error handling with workflow.GetLogger
- Implement idempotent activities
- Use appropriate retry policies
- Include OpenTelemetry tracing
- Follow Go best practices
"""
    
    return base_context


def build_user_prompt(request: ServiceGenerationRequest) -> str:
    """Build the user prompt with context."""
    prompt_parts = [request.prompt]
    
    if request.context:
        prompt_parts.append(f"\n\nAdditional context:\n{json.dumps(request.context, indent=2)}")
    
    if request.existing_service_id:
        prompt_parts.append(f"\n\nThis extends existing service: {request.existing_service_id}")
    
    return "\n".join(prompt_parts)


def select_local_model(mode: GenerationMode) -> str:
    """Select the best local model for the task."""
    if mode == GenerationMode.CODE_GENERATION:
        return "qwen2.5-coder-32b"
    return "llama-3.3-70b-instruct"


# Local Model Inference
@app.post("/api/v1/inference/local")
async def local_inference(request: InferenceRequest):
    """Direct inference on local models via vLLM."""
    with tracer.start_as_current_span("local_inference") as span:
        span.set_attribute("model", request.model)
        
        response = await local_vllm.complete(
            model=request.model,
            messages=request.messages,
            temperature=request.temperature,
            max_tokens=request.max_tokens,
        )
        
        return InferenceResponse(
            content=response.content,
            model=request.model,
            tokens_used=response.usage.total_tokens,
        )


# Workflow Enhancement
@app.post("/api/v1/enhance/workflow")
async def enhance_workflow(
    tenant_id: str,
    workflow_dsl: Dict[str, Any],
    enhancement_prompt: str,
):
    """
    Enhance an existing workflow based on natural language instructions.
    
    Examples:
    - "Add retry logic to all HTTP activities"
    - "Insert a notification step after payment processing"
    - "Add error handling for inventory check failures"
    """
    system_prompt = """You are a workflow enhancement expert.
Given an existing workflow DSL and an enhancement request, modify the workflow to implement the requested changes.
Preserve existing functionality while adding the new features.
Output the complete modified workflow DSL."""

    user_prompt = f"""
Existing workflow:
```json
{json.dumps(workflow_dsl, indent=2)}
```

Enhancement request: {enhancement_prompt}

Output the modified workflow DSL and explain the changes made.
"""

    response = await claude.complete(
        model="claude-sonnet-4-20250514",
        system=system_prompt,
        messages=[{"role": "user", "content": user_prompt}],
        temperature=0.3,
        max_tokens=8000,
    )
    
    # Parse and validate the enhanced workflow
    enhanced = parse_workflow_dsl(response.content)
    
    return {
        "enhanced_workflow": enhanced["workflow"],
        "changes": enhanced["explanation"],
    }


# Model Health Check
@app.get("/api/v1/models/health")
async def model_health():
    """Check health of all configured models."""
    health = {}
    
    # Check cloud providers
    for provider_name, provider in [("anthropic", claude), ("openai", openai)]:
        try:
            await provider.health_check()
            health[provider_name] = {"status": "healthy"}
        except Exception as e:
            health[provider_name] = {"status": "unhealthy", "error": str(e)}
    
    # Check local vLLM
    try:
        models = await local_vllm.list_models()
        health["local"] = {
            "status": "healthy",
            "models": models,
        }
    except Exception as e:
        health["local"] = {"status": "unhealthy", "error": str(e)}
    
    return health
```

### 5.3 Local Model Inference (vLLM)

```python
# ai-gateway/app/providers/local_vllm.py
import httpx
from typing import List, Dict, Any, Optional
from pydantic import BaseModel

from app.config import settings


class Message(BaseModel):
    role: str
    content: str


class CompletionRequest(BaseModel):
    model: str
    messages: List[Message]
    temperature: float = 0.7
    max_tokens: int = 4000
    top_p: float = 0.95
    stream: bool = False


class CompletionResponse(BaseModel):
    content: str
    model: str
    usage: Dict[str, int]


class VLLMProvider:
    def __init__(self, base_url: str):
        self.base_url = base_url
        self.client = httpx.AsyncClient(timeout=120.0)
    
    async def complete(
        self,
        model: str,
        messages: List[Dict[str, str]],
        temperature: float = 0.7,
        max_tokens: int = 4000,
        system: Optional[str] = None,
    ) -> CompletionResponse:
        """Execute completion against vLLM server."""
        
        # Prepend system message if provided
        if system:
            messages = [{"role": "system", "content": system}] + messages
        
        payload = {
            "model": model,
            "messages": messages,
            "temperature": temperature,
            "max_tokens": max_tokens,
            "top_p": 0.95,
        }
        
        response = await self.client.post(
            f"{self.base_url}/v1/chat/completions",
            json=payload,
        )
        response.raise_for_status()
        
        data = response.json()
        
        return CompletionResponse(
            content=data["choices"][0]["message"]["content"],
            model=model,
            usage={
                "prompt_tokens": data["usage"]["prompt_tokens"],
                "completion_tokens": data["usage"]["completion_tokens"],
                "total_tokens": data["usage"]["total_tokens"],
            },
        )
    
    async def list_models(self) -> List[str]:
        """List available models on vLLM server."""
        response = await self.client.get(f"{self.base_url}/v1/models")
        response.raise_for_status()
        
        data = response.json()
        return [m["id"] for m in data["data"]]
    
    async def health_check(self) -> bool:
        """Check if vLLM server is healthy."""
        response = await self.client.get(f"{self.base_url}/health")
        return response.status_code == 200


# Singleton instance
local_vllm = VLLMProvider(settings.VLLM_BASE_URL)
```

### 5.4 vLLM Kubernetes Deployment

```yaml
# k8s/ai-inference/vllm-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: vllm-inference
  namespace: omniroute-ai
spec:
  replicas: 1
  selector:
    matchLabels:
      app: vllm-inference
  template:
    metadata:
      labels:
        app: vllm-inference
    spec:
      nodeSelector:
        cloud.google.com/gke-accelerator: nvidia-tesla-a100
      tolerations:
      - key: "nvidia.com/gpu"
        operator: "Exists"
        effect: "NoSchedule"
      containers:
      - name: vllm
        image: vllm/vllm-openai:latest
        args:
        - "--model"
        - "meta-llama/Llama-3.3-70B-Instruct"
        - "--tensor-parallel-size"
        - "4"
        - "--max-model-len"
        - "32768"
        - "--gpu-memory-utilization"
        - "0.95"
        - "--enable-prefix-caching"
        ports:
        - containerPort: 8000
        resources:
          limits:
            nvidia.com/gpu: 4
            memory: "320Gi"
            cpu: "32"
          requests:
            nvidia.com/gpu: 4
            memory: "256Gi"
            cpu: "16"
        volumeMounts:
        - name: model-cache
          mountPath: /root/.cache/huggingface
        env:
        - name: HUGGING_FACE_HUB_TOKEN
          valueFrom:
            secretKeyRef:
              name: ai-secrets
              key: hf-token
        livenessProbe:
          httpGet:
            path: /health
            port: 8000
          initialDelaySeconds: 300
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /health
            port: 8000
          initialDelaySeconds: 120
          periodSeconds: 10
      volumes:
      - name: model-cache
        persistentVolumeClaim:
          claimName: model-cache-pvc
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: model-cache-pvc
  namespace: omniroute-ai
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 500Gi
  storageClassName: premium-rwo
---
apiVersion: v1
kind: Service
metadata:
  name: vllm-service
  namespace: omniroute-ai
spec:
  selector:
    app: vllm-inference
  ports:
  - port: 8000
    targetPort: 8000
  type: ClusterIP
```

---

## 6. Visual Service Designer (WYSIWYG)

### 6.1 React Flow Implementation

```typescript
// frontend/src/components/ServiceDesigner/ServiceDesigner.tsx
import React, { useCallback, useState, useRef } from 'react';
import ReactFlow, {
  Node,
  Edge,
  Connection,
  addEdge,
  useNodesState,
  useEdgesState,
  Controls,
  Background,
  MiniMap,
  Panel,
  NodeTypes,
  EdgeTypes,
} from 'reactflow';
import 'reactflow/dist/style.css';

import { ActivityNode } from './nodes/ActivityNode';
import { AIActionNode } from './nodes/AIActionNode';
import { N8NNode } from './nodes/N8NNode';
import { DecisionNode } from './nodes/DecisionNode';
import { ParallelNode } from './nodes/ParallelNode';
import { HumanTaskNode } from './nodes/HumanTaskNode';
import { NodePalette } from './NodePalette';
import { PropertyPanel } from './PropertyPanel';
import { AIAssistant } from './AIAssistant';
import { useServiceStore } from '@/stores/serviceStore';
import { compileWorkflow, validateWorkflow } from '@/lib/workflowCompiler';

const nodeTypes: NodeTypes = {
  activity: ActivityNode,
  ai_action: AIActionNode,
  n8n: N8NNode,
  decision: DecisionNode,
  parallel: ParallelNode,
  human_task: HumanTaskNode,
};

interface ServiceDesignerProps {
  serviceId?: string;
  tenantId: string;
  onSave: (workflow: WorkflowDSL) => Promise<void>;
}

export const ServiceDesigner: React.FC<ServiceDesignerProps> = ({
  serviceId,
  tenantId,
  onSave,
}) => {
  const reactFlowWrapper = useRef<HTMLDivElement>(null);
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const [selectedNode, setSelectedNode] = useState<Node | null>(null);
  const [validationErrors, setValidationErrors] = useState<string[]>([]);
  const [isAIAssistantOpen, setIsAIAssistantOpen] = useState(false);

  const { 
    availableActivities, 
    availableN8NWorkflows,
    aiModels,
  } = useServiceStore();

  // Handle new connections
  const onConnect = useCallback(
    (params: Connection) => {
      setEdges((eds) => addEdge({
        ...params,
        type: 'smoothstep',
        animated: true,
        style: { stroke: '#3b82f6', strokeWidth: 2 },
      }, eds));
    },
    [setEdges]
  );

  // Handle drag and drop from palette
  const onDragOver = useCallback((event: React.DragEvent) => {
    event.preventDefault();
    event.dataTransfer.dropEffect = 'move';
  }, []);

  const onDrop = useCallback(
    (event: React.DragEvent) => {
      event.preventDefault();

      if (!reactFlowWrapper.current) return;

      const reactFlowBounds = reactFlowWrapper.current.getBoundingClientRect();
      const nodeType = event.dataTransfer.getData('application/reactflow/type');
      const nodeData = JSON.parse(event.dataTransfer.getData('application/reactflow/data') || '{}');

      const position = {
        x: event.clientX - reactFlowBounds.left,
        y: event.clientY - reactFlowBounds.top,
      };

      const newNode: Node = {
        id: `${nodeType}_${Date.now()}`,
        type: nodeType,
        position,
        data: {
          label: nodeData.label || `New ${nodeType}`,
          config: nodeData.defaultConfig || {},
          ...nodeData,
        },
      };

      setNodes((nds) => [...nds, newNode]);
    },
    [setNodes]
  );

  // Handle node selection
  const onNodeClick = useCallback(
    (_: React.MouseEvent, node: Node) => {
      setSelectedNode(node);
    },
    []
  );

  // Update node configuration
  const onNodeConfigUpdate = useCallback(
    (nodeId: string, config: Record<string, unknown>) => {
      setNodes((nds) =>
        nds.map((node) =>
          node.id === nodeId
            ? { ...node, data: { ...node.data, config } }
            : node
        )
      );
    },
    [setNodes]
  );

  // Validate workflow
  const handleValidate = useCallback(async () => {
    const workflow = compileWorkflow(nodes, edges);
    const errors = await validateWorkflow(workflow);
    setValidationErrors(errors);
    return errors.length === 0;
  }, [nodes, edges]);

  // Save workflow
  const handleSave = useCallback(async () => {
    const isValid = await handleValidate();
    if (!isValid) {
      return;
    }

    const workflow = compileWorkflow(nodes, edges);
    await onSave(workflow);
  }, [nodes, edges, handleValidate, onSave]);

  // AI-powered workflow generation
  const handleAIGenerate = useCallback(
    async (prompt: string) => {
      const response = await fetch('/api/v1/generate/service', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          tenant_id: tenantId,
          mode: 'prompt_to_workflow',
          prompt,
          context: {
            existing_nodes: nodes.map(n => ({ type: n.type, label: n.data.label })),
            available_activities: availableActivities,
            available_n8n_workflows: availableN8NWorkflows,
          },
          use_local_model: false,
        }),
      });

      const data = await response.json();
      
      if (data.workflow_dsl) {
        // Convert DSL to React Flow nodes and edges
        const { newNodes, newEdges } = convertDSLToReactFlow(data.workflow_dsl);
        setNodes(newNodes);
        setEdges(newEdges);
      }

      return data;
    },
    [tenantId, nodes, availableActivities, availableN8NWorkflows, setNodes, setEdges]
  );

  // AI-powered workflow enhancement
  const handleAIEnhance = useCallback(
    async (enhancement: string) => {
      const currentWorkflow = compileWorkflow(nodes, edges);
      
      const response = await fetch('/api/v1/enhance/workflow', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          tenant_id: tenantId,
          workflow_dsl: currentWorkflow,
          enhancement_prompt: enhancement,
        }),
      });

      const data = await response.json();
      
      if (data.enhanced_workflow) {
        const { newNodes, newEdges } = convertDSLToReactFlow(data.enhanced_workflow);
        setNodes(newNodes);
        setEdges(newEdges);
      }

      return data;
    },
    [tenantId, nodes, edges, setNodes, setEdges]
  );

  return (
    <div className="h-screen flex">
      {/* Node Palette */}
      <NodePalette
        activities={availableActivities}
        n8nWorkflows={availableN8NWorkflows}
        aiModels={aiModels}
      />

      {/* Main Canvas */}
      <div className="flex-1" ref={reactFlowWrapper}>
        <ReactFlow
          nodes={nodes}
          edges={edges}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          onConnect={onConnect}
          onDrop={onDrop}
          onDragOver={onDragOver}
          onNodeClick={onNodeClick}
          nodeTypes={nodeTypes}
          fitView
          snapToGrid
          snapGrid={[15, 15]}
        >
          <Background gap={15} />
          <Controls />
          <MiniMap />
          
          {/* Top Toolbar */}
          <Panel position="top-center">
            <div className="flex gap-2 bg-white p-2 rounded-lg shadow-lg">
              <button
                onClick={handleValidate}
                className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
              >
                Validate
              </button>
              <button
                onClick={handleSave}
                className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
              >
                Save
              </button>
              <button
                onClick={() => setIsAIAssistantOpen(true)}
                className="px-4 py-2 bg-purple-500 text-white rounded hover:bg-purple-600"
              >
                🤖 AI Assistant
              </button>
            </div>
          </Panel>

          {/* Validation Errors */}
          {validationErrors.length > 0 && (
            <Panel position="bottom-center">
              <div className="bg-red-100 border border-red-400 text-red-700 p-4 rounded-lg max-w-xl">
                <h4 className="font-bold">Validation Errors:</h4>
                <ul className="list-disc list-inside">
                  {validationErrors.map((error, i) => (
                    <li key={i}>{error}</li>
                  ))}
                </ul>
              </div>
            </Panel>
          )}
        </ReactFlow>
      </div>

      {/* Property Panel */}
      {selectedNode && (
        <PropertyPanel
          node={selectedNode}
          onUpdate={onNodeConfigUpdate}
          onClose={() => setSelectedNode(null)}
          availableActivities={availableActivities}
          availableN8NWorkflows={availableN8NWorkflows}
          aiModels={aiModels}
        />
      )}

      {/* AI Assistant Modal */}
      {isAIAssistantOpen && (
        <AIAssistant
          onGenerate={handleAIGenerate}
          onEnhance={handleAIEnhance}
          onClose={() => setIsAIAssistantOpen(false)}
        />
      )}
    </div>
  );
};

// Helper function to convert DSL to React Flow format
function convertDSLToReactFlow(dsl: WorkflowDSL): {
  newNodes: Node[];
  newEdges: Edge[];
} {
  const newNodes: Node[] = dsl.nodes.map((node, index) => ({
    id: node.id,
    type: node.type,
    position: node.position || { x: 100 + (index % 3) * 250, y: 100 + Math.floor(index / 3) * 150 },
    data: {
      label: node.config.label || node.type,
      config: node.config,
    },
  }));

  const newEdges: Edge[] = dsl.edges.map((edge) => ({
    id: `${edge.source}-${edge.target}`,
    source: edge.source,
    target: edge.target,
    sourceHandle: edge.sourceHandle,
    targetHandle: edge.targetHandle,
    type: 'smoothstep',
    animated: true,
    style: { stroke: '#3b82f6', strokeWidth: 2 },
  }));

  return { newNodes, newEdges };
}
```

### 6.2 AI Assistant Component

```typescript
// frontend/src/components/ServiceDesigner/AIAssistant.tsx
import React, { useState, useCallback } from 'react';
import { motion, AnimatePresence } from 'framer-motion';

interface AIAssistantProps {
  onGenerate: (prompt: string) => Promise<GenerationResponse>;
  onEnhance: (enhancement: string) => Promise<EnhancementResponse>;
  onClose: () => void;
}

export const AIAssistant: React.FC<AIAssistantProps> = ({
  onGenerate,
  onEnhance,
  onClose,
}) => {
  const [mode, setMode] = useState<'generate' | 'enhance'>('generate');
  const [prompt, setPrompt] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [response, setResponse] = useState<any>(null);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = useCallback(async () => {
    if (!prompt.trim()) return;

    setIsLoading(true);
    setError(null);

    try {
      const result = mode === 'generate'
        ? await onGenerate(prompt)
        : await onEnhance(prompt);
      
      setResponse(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setIsLoading(false);
    }
  }, [mode, prompt, onGenerate, onEnhance]);

  const examplePrompts = mode === 'generate' ? [
    "Create a service that sends an SMS notification when an order is placed, waits 2 hours, and if not confirmed, escalates to a supervisor",
    "Build a workflow that fetches inventory from our ERP via n8n, uses AI to predict stock-outs, and creates purchase orders automatically",
    "Design a payment collection service that tries mobile money first, falls back to bank transfer, and records all attempts",
  ] : [
    "Add retry logic with exponential backoff to all HTTP activities",
    "Insert a notification step after the payment is confirmed",
    "Add error handling that creates a support ticket when inventory check fails",
    "Make the AI sentiment analysis step use a local model instead of cloud API",
  ];

  return (
    <motion.div
      initial={{ opacity: 0, y: 50 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: 50 }}
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
    >
      <div className="bg-white rounded-2xl shadow-2xl w-[800px] max-h-[80vh] overflow-hidden">
        {/* Header */}
        <div className="bg-gradient-to-r from-purple-600 to-blue-600 p-6 text-white">
          <div className="flex justify-between items-center">
            <div className="flex items-center gap-3">
              <span className="text-3xl">🤖</span>
              <div>
                <h2 className="text-xl font-bold">AI Service Assistant</h2>
                <p className="text-sm opacity-80">
                  Describe your service in plain English
                </p>
              </div>
            </div>
            <button
              onClick={onClose}
              className="p-2 hover:bg-white/20 rounded-full transition"
            >
              ✕
            </button>
          </div>

          {/* Mode Toggle */}
          <div className="flex gap-2 mt-4">
            <button
              onClick={() => setMode('generate')}
              className={`px-4 py-2 rounded-full transition ${
                mode === 'generate'
                  ? 'bg-white text-purple-600'
                  : 'bg-white/20 hover:bg-white/30'
              }`}
            >
              Generate New
            </button>
            <button
              onClick={() => setMode('enhance')}
              className={`px-4 py-2 rounded-full transition ${
                mode === 'enhance'
                  ? 'bg-white text-purple-600'
                  : 'bg-white/20 hover:bg-white/30'
              }`}
            >
              Enhance Existing
            </button>
          </div>
        </div>

        {/* Content */}
        <div className="p-6 overflow-y-auto max-h-[60vh]">
          {/* Input */}
          <div className="mb-6">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              {mode === 'generate'
                ? 'Describe the service you want to create:'
                : 'Describe the enhancement you want to make:'}
            </label>
            <textarea
              value={prompt}
              onChange={(e) => setPrompt(e.target.value)}
              placeholder={
                mode === 'generate'
                  ? "E.g., Create a service that monitors order status and sends notifications..."
                  : "E.g., Add retry logic to all payment activities..."
              }
              className="w-full h-32 p-4 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent resize-none"
            />
          </div>

          {/* Example Prompts */}
          <div className="mb-6">
            <p className="text-sm text-gray-500 mb-2">Try these examples:</p>
            <div className="flex flex-wrap gap-2">
              {examplePrompts.map((example, i) => (
                <button
                  key={i}
                  onClick={() => setPrompt(example)}
                  className="px-3 py-1 text-sm bg-gray-100 hover:bg-gray-200 rounded-full text-gray-700 truncate max-w-xs"
                  title={example}
                >
                  {example.length > 50 ? example.substring(0, 50) + '...' : example}
                </button>
              ))}
            </div>
          </div>

          {/* Response */}
          {response && (
            <div className="mb-6 p-4 bg-green-50 border border-green-200 rounded-lg">
              <h4 className="font-medium text-green-800 mb-2">
                ✅ {mode === 'generate' ? 'Service Generated!' : 'Workflow Enhanced!'}
              </h4>
              {response.explanation && (
                <p className="text-sm text-gray-700 mb-2">{response.explanation}</p>
              )}
              {response.suggestions && (
                <div className="mt-2">
                  <p className="text-sm font-medium text-gray-600">Suggestions:</p>
                  <ul className="list-disc list-inside text-sm text-gray-600">
                    {response.suggestions.map((s: string, i: number) => (
                      <li key={i}>{s}</li>
                    ))}
                  </ul>
                </div>
              )}
              <div className="mt-2 text-xs text-gray-500">
                Model: {response.model_used} | Tokens: {response.tokens_used} | Latency: {response.latency_ms}ms
              </div>
            </div>
          )}

          {/* Error */}
          {error && (
            <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
              <p className="text-red-700">{error}</p>
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="border-t p-4 flex justify-end gap-3">
          <button
            onClick={onClose}
            className="px-6 py-2 text-gray-600 hover:bg-gray-100 rounded-lg transition"
          >
            Cancel
          </button>
          <button
            onClick={handleSubmit}
            disabled={isLoading || !prompt.trim()}
            className={`px-6 py-2 bg-purple-600 text-white rounded-lg transition ${
              isLoading || !prompt.trim()
                ? 'opacity-50 cursor-not-allowed'
                : 'hover:bg-purple-700'
            }`}
          >
            {isLoading ? (
              <span className="flex items-center gap-2">
                <span className="animate-spin">⏳</span>
                Generating...
              </span>
            ) : mode === 'generate' ? (
              'Generate Service'
            ) : (
              'Enhance Workflow'
            )}
          </button>
        </div>
      </div>
    </motion.div>
  );
};
```

---

## 7. GCP Infrastructure & Deployment

### 7.1 GKE Cluster Architecture

```yaml
# terraform/gcp/gke.tf
resource "google_container_cluster" "omniroute_sce" {
  name     = "omniroute-sce-cluster"
  location = "us-central1"
  
  # Enable Autopilot for serverless K8s
  enable_autopilot = true
  
  # Network configuration
  network    = google_compute_network.sce_network.name
  subnetwork = google_compute_subnetwork.sce_subnet.name
  
  # IP allocation policy
  ip_allocation_policy {
    cluster_secondary_range_name  = "pods"
    services_secondary_range_name = "services"
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
  
  # Maintenance window
  maintenance_policy {
    daily_maintenance_window {
      start_time = "03:00"
    }
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
}

# Node pool for GPU workloads
resource "google_container_node_pool" "gpu_pool" {
  name       = "gpu-pool"
  location   = "us-central1"
  cluster    = google_container_cluster.omniroute_sce.name
  
  node_count = 0  # Start with 0, autoscale up
  
  autoscaling {
    min_node_count = 0
    max_node_count = 4
  }
  
  node_config {
    machine_type = "a2-highgpu-4g"  # 4x A100 GPUs
    
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
  }
}

# Cloud SQL for PostgreSQL (Temporal & n8n)
resource "google_sql_database_instance" "sce_postgres" {
  name             = "omniroute-sce-postgres"
  database_version = "POSTGRES_15"
  region           = "us-central1"
  
  settings {
    tier = "db-custom-8-32768"  # 8 vCPUs, 32GB RAM
    
    availability_type = "REGIONAL"
    
    backup_configuration {
      enabled                        = true
      point_in_time_recovery_enabled = true
      backup_retention_settings {
        retained_backups = 30
      }
    }
    
    ip_configuration {
      ipv4_enabled    = false
      private_network = google_compute_network.sce_network.id
    }
    
    database_flags {
      name  = "max_connections"
      value = "500"
    }
    
    database_flags {
      name  = "work_mem"
      value = "64MB"
    }
    
    insights_config {
      query_insights_enabled  = true
      query_string_length     = 2048
      record_application_tags = true
      record_client_address   = true
    }
  }
}

# Memorystore Redis
resource "google_redis_instance" "sce_redis" {
  name           = "omniroute-sce-redis"
  tier           = "STANDARD_HA"
  memory_size_gb = 16
  region         = "us-central1"
  
  redis_version     = "REDIS_7_0"
  display_name      = "OmniRoute SCE Redis"
  
  authorized_network = google_compute_network.sce_network.id
  
  persistence_config {
    persistence_mode    = "RDB"
    rdb_snapshot_period = "ONE_HOUR"
  }
  
  maintenance_policy {
    weekly_maintenance_window {
      day = "SUNDAY"
      start_time {
        hours   = 3
        minutes = 0
      }
    }
  }
}

# Cloud Armor for DDoS protection
resource "google_compute_security_policy" "sce_policy" {
  name = "omniroute-sce-security-policy"
  
  # Default rule
  rule {
    action   = "allow"
    priority = "2147483647"
    match {
      versioned_expr = "SRC_IPS_V1"
      config {
        src_ip_ranges = ["*"]
      }
    }
    description = "Default allow rule"
  }
  
  # Block known bad IPs
  rule {
    action   = "deny(403)"
    priority = "1000"
    match {
      expr {
        expression = "evaluatePreconfiguredExpr('xss-stable')"
      }
    }
    description = "XSS protection"
  }
  
  # Rate limiting
  rule {
    action   = "rate_based_ban"
    priority = "2000"
    rate_limit_options {
      conform_action = "allow"
      exceed_action  = "deny(429)"
      rate_limit_threshold {
        count        = 1000
        interval_sec = 60
      }
      ban_duration_sec = 600
    }
    match {
      versioned_expr = "SRC_IPS_V1"
      config {
        src_ip_ranges = ["*"]
      }
    }
    description = "Rate limiting rule"
  }
}
```

### 7.2 Helm Chart for SCE Deployment

```yaml
# helm/omniroute-sce/values.yaml
global:
  namespace: omniroute-sce
  image:
    registry: gcr.io/omniroute-platform
    pullPolicy: IfNotPresent
  
  # Environment
  environment: production
  
  # Observability
  otel:
    enabled: true
    endpoint: "http://otel-collector.observability:4317"
  
  # Database
  postgres:
    host: postgres-sce.omniroute-sce.svc.cluster.local
    port: 5432
    database: sce
    sslMode: require
  
  # Redis
  redis:
    host: redis-sce.omniroute-sce.svc.cluster.local
    port: 6379
    tls: true

# Temporal Server
temporal:
  enabled: true
  replicaCount: 3
  
  server:
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
      persistence:
        default:
          sql:
            driver: postgres
            host: ${global.postgres.host}
            port: ${global.postgres.port}
            database: temporal
            user: temporal
            existingSecret: temporal-db-credentials
        visibility:
          sql:
            driver: postgres
            host: ${global.postgres.host}
            port: ${global.postgres.port}
            database: temporal_visibility
            user: temporal
            existingSecret: temporal-db-credentials

# Temporal Workers
temporalWorkers:
  core:
    enabled: true
    replicaCount: 5
    image: ${global.image.registry}/sce-worker-core
    tag: latest
    
    taskQueues:
      - omniroute-core
    
    resources:
      requests:
        memory: "512Mi"
        cpu: "250m"
      limits:
        memory: "2Gi"
        cpu: "1000m"
    
    autoscaling:
      enabled: true
      minReplicas: 3
      maxReplicas: 20
      targetCPUUtilization: 70
  
  integration:
    enabled: true
    replicaCount: 3
    image: ${global.image.registry}/sce-worker-integration
    tag: latest
    
    taskQueues:
      - omniroute-integration
    
    resources:
      requests:
        memory: "256Mi"
        cpu: "100m"
      limits:
        memory: "1Gi"
        cpu: "500m"
  
  ai:
    enabled: true
    replicaCount: 2
    image: ${global.image.registry}/sce-worker-ai
    tag: latest
    
    taskQueues:
      - omniroute-ai
    
    resources:
      requests:
        memory: "1Gi"
        cpu: "500m"
      limits:
        memory: "4Gi"
        cpu: "2000m"

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
  
  worker:
    replicaCount: 5
    
    autoscaling:
      enabled: true
      minReplicas: 3
      maxReplicas: 15
      targetCPUUtilization: 70

# AI Gateway
aiGateway:
  enabled: true
  replicaCount: 3
  image: ${global.image.registry}/sce-ai-gateway
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

# vLLM Inference
vllm:
  enabled: true
  replicaCount: 1
  image: vllm/vllm-openai
  tag: latest
  
  model: meta-llama/Llama-3.3-70B-Instruct
  tensorParallelSize: 4
  maxModelLen: 32768
  
  resources:
    limits:
      nvidia.com/gpu: 4
      memory: "320Gi"
      cpu: "32"

# Service Registry API
serviceRegistry:
  enabled: true
  replicaCount: 3
  image: ${global.image.registry}/sce-service-registry
  tag: latest
  
  resources:
    requests:
      memory: "256Mi"
      cpu: "100m"
    limits:
      memory: "1Gi"
      cpu: "500m"

# Workflow Compiler
workflowCompiler:
  enabled: true
  replicaCount: 3
  image: ${global.image.registry}/sce-workflow-compiler
  tag: latest
  
  resources:
    requests:
      memory: "256Mi"
      cpu: "100m"
    limits:
      memory: "512Mi"
      cpu: "250m"

# Frontend
frontend:
  enabled: true
  replicaCount: 2
  image: ${global.image.registry}/sce-frontend
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
    kubernetes.io/ingress.class: gce
    cloud.google.com/backend-config: '{"default": "sce-backend-config"}'
  
  hosts:
    - host: sce.omniroute.io
      paths:
        - path: /
          pathType: Prefix
          service:
            name: sce-frontend
            port: 80
        - path: /api
          pathType: Prefix
          service:
            name: sce-api-gateway
            port: 80
        - path: /n8n
          pathType: Prefix
          service:
            name: n8n-main
            port: 5678
```

---

## 8. XP Practices Implementation

### 8.1 Test-Driven Development (TDD)

```go
// internal/domain/service_definition_test.go
package domain_test

import (
    "testing"
    "time"
    
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    
    "github.com/omniroute/sce/internal/domain"
)

func TestServiceDefinition_Publish(t *testing.T) {
    t.Parallel()
    
    tests := []struct {
        name        string
        setup       func() *domain.ServiceDefinition
        wantErr     error
        wantStatus  domain.ServiceStatus
        wantEvents  int
    }{
        {
            name: "successfully publishes draft service with workflow",
            setup: func() *domain.ServiceDefinition {
                return &domain.ServiceDefinition{
                    ID:       domain.ServiceID(uuid.New()),
                    TenantID: domain.TenantID(uuid.New()),
                    Name:     "Test Service",
                    Status:   domain.ServiceStatusDraft,
                    Workflow: domain.WorkflowGraph{
                        Nodes: []domain.WorkflowNode{
                            {ID: "1", Type: domain.NodeTypeActivity},
                        },
                    },
                }
            },
            wantErr:    nil,
            wantStatus: domain.ServiceStatusPublished,
            wantEvents: 1,
        },
        {
            name: "fails to publish archived service",
            setup: func() *domain.ServiceDefinition {
                return &domain.ServiceDefinition{
                    ID:     domain.ServiceID(uuid.New()),
                    Status: domain.ServiceStatusArchived,
                    Workflow: domain.WorkflowGraph{
                        Nodes: []domain.WorkflowNode{{ID: "1"}},
                    },
                }
            },
            wantErr:    domain.ErrCannotPublishArchived,
            wantStatus: domain.ServiceStatusArchived,
            wantEvents: 0,
        },
        {
            name: "fails to publish service without workflow",
            setup: func() *domain.ServiceDefinition {
                return &domain.ServiceDefinition{
                    ID:       domain.ServiceID(uuid.New()),
                    Status:   domain.ServiceStatusDraft,
                    Workflow: domain.WorkflowGraph{},
                }
            },
            wantErr:    domain.ErrEmptyWorkflow,
            wantStatus: domain.ServiceStatusDraft,
            wantEvents: 0,
        },
    }
    
    for _, tt := range tests {
        tt := tt
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            
            svc := tt.setup()
            
            err := svc.Publish()
            
            if tt.wantErr != nil {
                require.Error(t, err)
                assert.Equal(t, tt.wantErr, err)
            } else {
                require.NoError(t, err)
            }
            
            assert.Equal(t, tt.wantStatus, svc.Status)
            assert.Len(t, svc.PullEvents(), tt.wantEvents)
        })
    }
}

func TestServiceDefinition_AddVersion(t *testing.T) {
    t.Parallel()
    
    t.Run("creates new version with incremented number", func(t *testing.T) {
        svc := &domain.ServiceDefinition{
            ID:       domain.ServiceID(uuid.New()),
            Versions: []domain.ServiceVersion{},
        }
        
        newWorkflow := domain.WorkflowGraph{
            Nodes: []domain.WorkflowNode{{ID: "new"}},
        }
        
        versionID, err := svc.AddVersion(newWorkflow, "Initial release")
        
        require.NoError(t, err)
        assert.NotEqual(t, uuid.Nil, uuid.UUID(versionID))
        assert.Len(t, svc.Versions, 1)
        assert.Equal(t, versionID, svc.ActiveVersion)
        assert.Equal(t, 1, svc.Versions[0].VersionNumber)
        
        // Add another version
        versionID2, err := svc.AddVersion(newWorkflow, "Bug fix")
        require.NoError(t, err)
        assert.Len(t, svc.Versions, 2)
        assert.Equal(t, versionID2, svc.ActiveVersion)
        assert.Equal(t, 2, svc.Versions[1].VersionNumber)
    })
}

func TestServiceName_Validate(t *testing.T) {
    t.Parallel()
    
    tests := []struct {
        name    string
        input   domain.ServiceName
        wantErr bool
    }{
        {"valid name", "My Service", false},
        {"minimum length", "abc", false},
        {"maximum length", string(make([]byte, 100)), false},
        {"too short", "ab", true},
        {"too long", string(make([]byte, 101)), true},
        {"empty", "", true},
    }
    
    for _, tt := range tests {
        tt := tt
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            
            err := tt.input.Validate()
            
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### 8.2 Continuous Integration Pipeline

```yaml
# .github/workflows/ci.yaml
name: CI Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

env:
  GO_VERSION: '1.23'
  RUST_VERSION: '1.75'
  NODE_VERSION: '20'
  PYTHON_VERSION: '3.12'

jobs:
  # Go Services
  go-test:
    name: Go Tests
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_PASSWORD: test
          POSTGRES_DB: test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
      redis:
        image: redis:7
        ports:
          - 6379:6379
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
      
      - name: Cache Go modules
        uses: actions/cache@v4
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
      
      - name: Install dependencies
        run: go mod download
      
      - name: Run tests
        run: |
          go test -v -race -coverprofile=coverage.out ./...
          go tool cover -html=coverage.out -o coverage.html
      
      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: ./coverage.out
      
      - name: Run linter
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest

  # Rust Workflow Compiler
  rust-test:
    name: Rust Tests
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Rust
        uses: dtolnay/rust-toolchain@stable
        with:
          toolchain: ${{ env.RUST_VERSION }}
          components: clippy, rustfmt
      
      - name: Cache Cargo
        uses: Swatinem/rust-cache@v2
      
      - name: Check formatting
        run: cargo fmt --all -- --check
      
      - name: Run Clippy
        run: cargo clippy --all-targets --all-features -- -D warnings
      
      - name: Run tests
        run: cargo test --all-features

  # Python AI Gateway
  python-test:
    name: Python Tests
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Python
        uses: actions/setup-python@v5
        with:
          python-version: ${{ env.PYTHON_VERSION }}
      
      - name: Install Poetry
        run: pip install poetry
      
      - name: Install dependencies
        run: |
          cd ai-gateway
          poetry install
      
      - name: Run tests
        run: |
          cd ai-gateway
          poetry run pytest --cov=app --cov-report=xml
      
      - name: Run linter
        run: |
          cd ai-gateway
          poetry run ruff check .
          poetry run mypy app

  # Frontend
  frontend-test:
    name: Frontend Tests
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Node.js
        uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
          cache: 'pnpm'
      
      - name: Install pnpm
        run: npm install -g pnpm
      
      - name: Install dependencies
        run: |
          cd frontend
          pnpm install
      
      - name: Run type check
        run: |
          cd frontend
          pnpm type-check
      
      - name: Run linter
        run: |
          cd frontend
          pnpm lint
      
      - name: Run tests
        run: |
          cd frontend
          pnpm test:coverage
      
      - name: Build
        run: |
          cd frontend
          pnpm build

  # Integration Tests
  integration-test:
    name: Integration Tests
    runs-on: ubuntu-latest
    needs: [go-test, rust-test, python-test]
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Docker Compose
        run: docker compose -f docker-compose.test.yaml up -d
      
      - name: Wait for services
        run: |
          sleep 30
          docker compose -f docker-compose.test.yaml ps
      
      - name: Run integration tests
        run: |
          go test -v -tags=integration ./tests/integration/...
      
      - name: Cleanup
        run: docker compose -f docker-compose.test.yaml down -v

  # Build and Push Images
  build:
    name: Build Images
    runs-on: ubuntu-latest
    needs: [integration-test, frontend-test]
    if: github.ref == 'refs/heads/main'
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3
      
      - name: Login to GCR
        uses: docker/login-action@v3
        with:
          registry: gcr.io
          username: _json_key
          password: ${{ secrets.GCP_SA_KEY }}
      
      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: |
            gcr.io/omniroute-platform/sce-service-registry:${{ github.sha }}
            gcr.io/omniroute-platform/sce-service-registry:latest
          cache-from: type=gha
          cache-to: type=gha,mode=max
```

---

## 9. Implementation Prompts

### Prompt 1: Domain Layer Implementation

```
You are implementing the Domain Layer for the OmniRoute Service Creation Environment (SCE) following strict DDD principles.

**Context:**
- The SCE allows non-technical users to create business services via visual design and AI prompts
- Services are executed via Temporal workflows with n8n integrations and AI capabilities
- Multi-tenant architecture with strict data isolation

**Requirements:**

1. **Bounded Context: Service Definition**
   - Implement the `ServiceDefinition` aggregate root with:
     - Identity: ID, TenantID
     - State: Name, Description, Category, Status, Workflow
     - Behavior: Publish(), Deprecate(), Archive(), AddVersion()
     - Invariants: Cannot publish without workflow, cannot archive if active instances
   
2. **Value Objects:**
   - `ServiceName`: Validated string (3-100 chars, no special chars)
   - `WorkflowGraph`: Immutable graph of nodes and edges
   - `WorkflowNode`: Node with type, config, position
   - `ActivityReference`: Reference to a Temporal activity

3. **Domain Events:**
   - `ServiceCreated`
   - `ServicePublished`
   - `ServiceDeprecated`
   - `VersionReleased`

4. **Repository Interface:**
   - `ServiceDefinitionRepository` with methods for CRUD and queries
   - Must support optimistic locking
   - Must emit domain events

**Constraints:**
- Use Go 1.23+
- All domain logic must be pure (no I/O)
- 100% test coverage for domain methods
- Use table-driven tests
- Follow Ubiquitous Language from the domain model

**Output:**
Generate the complete domain layer package including:
- types.go (all value objects and entities)
- service_definition.go (aggregate root)
- events.go (domain events)
- errors.go (domain errors)
- repository.go (repository interface)
- *_test.go files for all packages

Include detailed comments explaining DDD concepts.
```

### Prompt 2: Temporal Workflow Implementation

```
You are implementing the Temporal Workflow Engine for the OmniRoute SCE.

**Context:**
The workflow engine executes user-defined services as Temporal workflows. Each workflow is compiled from a visual DSL into executable Go code.

**Requirements:**

1. **Core Workflows:**
   - `UserDefinedServiceWorkflow`: Executes any user-created service
   - Must handle dynamic node execution based on DSL
   - Support all node types: Activity, AI Action, n8n, Decision, Parallel, Wait, Human Task

2. **Activities:**
   - `CoreActivities`: Database, HTTP, Notification
   - `IntegrationActivities`: n8n execution, webhook handling
   - `AIActivities`: LLM calls, local model inference

3. **Features:**
   - Retry policies per activity type
   - Saga pattern for compensations
   - Signals for human-in-the-loop
   - Queries for workflow state
   - Child workflows for subflows

4. **Performance:**
   - Connection pooling for external calls
   - Circuit breakers for unreliable services
   - Metrics and tracing via OpenTelemetry

**Constraints:**
- Use Temporal Go SDK v1.24+
- Activities must be idempotent
- Use heartbeats for long-running activities
- Implement proper error handling with retries

**Output:**
Generate:
- workflows/user_service.go (main workflow)
- workflows/common.go (shared utilities)
- activities/core.go
- activities/integration.go
- activities/ai.go
- workers/main.go (worker setup)
- All with comprehensive tests
```

### Prompt 3: AI Service Generation

```
You are implementing the AI Service Generation Engine for OmniRoute SCE.

**Context:**
Users can describe services in natural language, and the AI engine generates complete service definitions with workflows.

**Requirements:**

1. **Generation Modes:**
   - `PromptToService`: Full service from description
   - `PromptToWorkflow`: Workflow from existing service context
   - `CodeGeneration`: Temporal Go code from DSL
   - `Enhancement`: Modify existing workflow via prompt

2. **Model Support:**
   - Cloud LLMs: Claude (Anthropic), GPT-4o (OpenAI)
   - Local Models: Llama 3.3, Mistral, Qwen (via vLLM)
   - Automatic fallback and load balancing

3. **Prompt Engineering:**
   - System prompts per generation mode
   - Context injection (available activities, n8n workflows)
   - Output parsing and validation
   - Structured output via JSON mode

4. **Features:**
   - Rate limiting per tenant
   - Cost tracking per model
   - Response caching
   - Streaming support

**Constraints:**
- Use FastAPI with async
- Implement circuit breakers
- Add OpenTelemetry tracing
- Type hints throughout

**Output:**
Generate Python packages:
- app/main.py (FastAPI app)
- app/providers/claude.py
- app/providers/openai.py
- app/providers/local_vllm.py
- app/services/generator.py
- app/prompts/templates.py
- tests/ with full coverage
```

### Prompt 4: Visual Service Designer

```
You are implementing the Visual Service Designer for OmniRoute SCE using React and React Flow.

**Context:**
Non-technical users design services by dragging nodes onto a canvas and connecting them.

**Requirements:**

1. **Canvas Features:**
   - Drag-and-drop nodes from palette
   - Connect nodes with edges
   - Undo/redo support
   - Zoom, pan, minimap
   - Grid snapping

2. **Node Types:**
   - Activity: Configure Temporal activity
   - AI Action: Configure LLM/local model call
   - n8n: Select and configure n8n workflow
   - Decision: Conditional branching with expression editor
   - Parallel: Fork/join for concurrent execution
   - Wait: Timer or signal wait
   - Human Task: Approval/input collection

3. **Property Panel:**
   - Dynamic form based on node type
   - JSON schema validation
   - Expression editor with autocomplete
   - Preview mode

4. **AI Assistant:**
   - Chat interface for natural language generation
   - Enhancement suggestions
   - Error explanation

5. **Export:**
   - Compile to workflow DSL
   - Validate before save
   - Version comparison

**Constraints:**
- React 18+ with TypeScript
- React Flow for canvas
- Zustand for state
- Tailwind CSS for styling
- Vitest for testing

**Output:**
Generate:
- components/ServiceDesigner/
- components/nodes/ (all node components)
- components/PropertyPanel/
- components/AIAssistant/
- stores/serviceStore.ts
- lib/workflowCompiler.ts
- Full test coverage
```

### Prompt 5: GCP Infrastructure

```
You are implementing the GCP infrastructure for OmniRoute SCE using Terraform and Helm.

**Context:**
Deploy the SCE platform on GKE with high availability, auto-scaling, and security.

**Requirements:**

1. **GKE Cluster:**
   - Autopilot mode for serverless K8s
   - Multi-zone deployment
   - Workload Identity for IAM
   - Binary Authorization

2. **GPU Node Pool:**
   - A100 GPUs for AI inference
   - Auto-scaling 0-4 nodes
   - Preemptible for cost savings

3. **Databases:**
   - Cloud SQL PostgreSQL (HA)
   - Memorystore Redis (HA)
   - Cloud Storage for model cache

4. **Networking:**
   - VPC with private subnets
   - Cloud NAT for egress
   - Cloud Armor for DDoS
   - Cloud CDN for static assets

5. **Security:**
   - Secret Manager for credentials
   - KMS for encryption
   - IAM least privilege

6. **Observability:**
   - Cloud Monitoring
   - Cloud Logging
   - Cloud Trace

**Constraints:**
- Terraform 1.6+
- Helm 3.13+
- Follow GCP best practices
- Cost optimization

**Output:**
Generate:
- terraform/gcp/ (all .tf files)
- helm/omniroute-sce/ (chart)
- scripts/deploy.sh
- Documentation for deployment
```

---

## 10. Performance Optimization

### 10.1 Performance Targets

| Metric | Target | Measurement |
|--------|--------|-------------|
| Service Creation API | < 100ms p99 | Create service definition |
| Workflow Compilation | < 500ms p99 | DSL to Go code |
| AI Generation | < 5s p99 | Prompt to service |
| Local Inference | < 2s p99 | Llama 70B completion |
| n8n Execution | < 10s p99 | Average workflow |
| Temporal Workflow Start | < 50ms p99 | Start latency |
| Canvas Render | 60 FPS | 100+ nodes |

### 10.2 Optimization Strategies

1. **Workflow Compiler (Rust)**
   - Incremental compilation with caching
   - Parallel code generation
   - Pre-compiled templates

2. **AI Inference**
   - vLLM with PagedAttention
   - Continuous batching
   - Speculative decoding
   - Model quantization (INT8/FP8)

3. **Frontend**
   - React Flow virtualization
   - Web Workers for validation
   - IndexedDB for offline support

4. **Database**
   - Connection pooling (PgBouncer)
   - Read replicas for queries
   - Partitioning for large tables

5. **Caching**
   - Redis for hot data
   - CDN for static assets
   - Browser caching

---

## Appendix A: API Reference

[Full OpenAPI specification available in /api/openapi.yaml]

## Appendix B: Event Schema Registry

[Avro schemas available in /schemas/events/]

## Appendix C: Monitoring Dashboards

[Grafana dashboards available in /monitoring/dashboards/]

---

**Document Control:**
- Author: OmniRoute Architecture Team
- Reviewers: Engineering Leadership
- Last Updated: January 2026
- Next Review: April 2026
