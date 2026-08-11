#![allow(dead_code)]

use crate::scheme::{
    builtins, parser::{Error, Literal, Scope, Token, lex}, utils::describe,
};
use log::{debug, error, trace, warn};
use std::collections::HashMap;

type BuiltinFn = fn(&mut Engine) -> Result<(), String>;

#[derive(Clone, Copy, Debug, PartialEq)]
#[allow(unpredictable_function_pointer_comparisons)]
pub struct NativeFn {
    argc: Option<u32>,
    call: BuiltinFn,
}

impl NativeFn {
    pub fn check_args(&self, scope: &Scope) -> Result<(), Error> {
        let args = scope.args.len();

        if let Some(a) = self.argc
        && (args > a as usize || args < a as usize)
        {
            {
                return Err(Error::new(
                    scope.line,
                    format!(
                        "Too {} arguments\n(Expected {a}, got {args})",
                        if args < a as usize { "few" } else { "many" },
                    ),
                ));
            }
        }

        Ok(())
    }
}

#[derive(Clone, Debug, PartialEq)]
pub enum Value {
    /// Atom is a single, non-quoted word used as enumeration
    Atom(String),
    Str(String),
    Int(i64),
    Float(f64),
    Lit(Literal),
    Builtin(NativeFn),
    Bool(bool),
}

#[derive(Clone, Debug, Default)]
pub struct EngineConfig {
    pub display_flushes: bool,
}

#[derive(Default)]
pub struct Engine {
    global: HashMap<String, Value>,
    pub stack: Vec<Value>,

    pub config: EngineConfig
}

impl Engine {
    pub fn new() -> Self {
        let mut e = Engine::default();

        e.reg_builtin("display", Some(1), builtins::display);
        e.reg_builtin("newline", Some(0), builtins::newline);
        
        e.reg_builtin("boolean?", Some(1), builtins::is_bool);
        e.reg_builtin("number?", Some(1), builtins::is_number);
        e.reg_builtin("integer?", Some(1), builtins::is_int);
        e.reg_builtin("eqv?", Some(2), builtins::eqv);
        
        e.reg_builtin("+", Some(2), builtins::add);

        e
    }

    pub fn reg_builtin(&mut self, name: &'static str, argc: Option<u32>, fun: BuiltinFn) {
        self.reg_global(
            name.to_string(),
            Value::Builtin(NativeFn {
                call: fun,
                argc: argc,
            }),
        );
    }

    pub fn reg_global(&mut self, name: String, val: Value) {
        if self.global.contains_key(&name) {
            warn!("Global value '{name}' will be overwritten");
        }

        self.global.insert(name, val);
    }

    pub fn call_fn(&mut self, scope: &Scope) -> Result<(), Error> {
        if !self.global.contains_key(&scope.name) {
            Err(Error::new(
                scope.line,
                format!("function '{}' not found", scope.name),
            ))
        } else {
            let f = if let Value::Builtin(fun) = *self.global.get(&scope.name.clone()).unwrap() {
                Ok(fun)
            } else {
                Err(Error::new(
                    scope.line,
                    format!("{} is not a function", scope.name),
                ))
            }?;

            f.check_args(scope)?;
            self.populate_stack(&scope.args)?;

            match (f.call)(self) {
                Err(e) => Err(Error::new(scope.line, format!("{} failed: {}", scope.name, e))),
                _ => Ok(()),
            }
        }
    }

    pub fn run_str(&mut self, src: String, str: String) {
        let tokens = match lex(str) {
            Ok(toks) => toks,
            Err(e) => {
                error!("{src}:{e}");
                return;
            }
        };

        debug!("GOT: {}", describe(tokens.clone(), 0));

        for scope in tokens {
            match self.call_fn(&scope) {
                Err(e) => error!("{src}:{e}"),
                _ => {}
            }
        }

        if let Some(v) = self.stack.pop() {
            trace!("-> {v}");
        }
    }
    
    fn populate_stack(&mut self, args: &Vec<Token>) -> Result<(), Error> {
        for arg in args {
            match arg.kind.clone() {
                super::parser::TokenKind::Ident(i) => self.stack.push(Value::Str(i)),
                super::parser::TokenKind::Literal(li) => self.stack.push(Value::Lit(li)),
                // The function should push it's return into the stack, setting it as the next arg
                super::parser::TokenKind::Scope(scope) => self.call_fn(&scope)?, 
            }
        }

        Ok(())
    }
}
