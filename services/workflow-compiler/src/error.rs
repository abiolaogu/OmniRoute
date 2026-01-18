//! Error types for Workflow Compiler

use thiserror::Error;

#[derive(Error, Debug)]
pub enum CompilerError {
    #[error("Validation error: {0}")]
    ValidationError(String),
    
    #[error("Parse error: {0}")]
    ParseError(String),
    
    #[error("Cycle detected in workflow graph")]
    CycleDetected,
    
    #[error("Code generation error: {0}")]
    CodeGenError(String),
    
    #[error("Template error: {0}")]
    TemplateError(#[from] handlebars::TemplateError),
    
    #[error("IO error: {0}")]
    IoError(#[from] std::io::Error),
}
