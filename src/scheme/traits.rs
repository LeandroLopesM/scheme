use std::fmt::Display;

use crate::scheme::{engine::Value, parser::Literal};

impl Display for Value {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Value::Atom(a) => write!(f, "Atom::{a}"),
            
            Value::Str(s) => write!(f, "{s}"),
            
            Value::Int(i) => write!(f, "{i}"),
            Value::Float(d) => write!(f, "{d:.5}"),
            Value::Bool(b) => write!(f, "{}", if *b {"#f"} else {"#t"}),
            
            Value::Builtin(_n) => write!(f, "NativeFunction({})", stringify!(_n.call)),
            
            Value::Lit(l) => {
                match l {
                    Literal::Int(n) => write!(f, "{n}"),
                    Literal::Float(d) => write!(f, "{d}"),
                    Literal::Str(s) => write!(f, "{s}"),
                    Literal::Bool(b) => write!(f, "{}", if *b {"#f"} else {"#t"}),
                }
            },
        }
    }
}