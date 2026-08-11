use std::{fmt::Display, ops::Add};

use crate::scheme::{engine::{NativeFn}, parser::Literal};

#[derive(Clone, Debug, PartialEq, PartialOrd)]
pub enum Number {
    Int(i64),       // | | | | Integer
    Rat(f64),       // | | | Rational
    Real(f64),      // | | Real
    Comp(f64),      // | Complex
                    // Number (wildcard)
}

impl Display for Number {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Number::Int(v) => write!(f, "{v}"),
            Number::Rat(v) => write!(f, "{v}"),
            Number::Real(v) => write!(f, "{v}"),
            Number::Comp(v) => write!(f, "{v}"),
        }
    }
}

impl Add<Number> for Number {
    type Output = Self;

    fn add(self, rhs: Number) -> Self::Output {
        let lhs = match self {
            Number::Int(v) => v as f64,
            Number::Rat(v) => v,
            Number::Real(v) => v,
            Number::Comp(v) => v,
        };
        
        let other = match self {
            Number::Int(v) => v as f64,
            Number::Rat(v) => v,
            Number::Real(v) => v,
            Number::Comp(v) => v,
        };


        let res = lhs + other;

        let matcher = if self < rhs {self} else {rhs};

        match matcher { // Output datatype is whichever one was larger
            Number::Int(_) => Number::Int(res as i64),
            Number::Rat(_) => Number::Rat(res),
            Number::Real(_) => Number::Real(res),
            Number::Comp(_) => Number::Comp(res),
        }
    }
}

#[derive(Clone, Debug, PartialEq)]
pub enum Value {
    /// Atom is a single, non-quoted word used as enumeration
    Atom(String),
    Number(Number),

    Str(String),
    Lit(Literal),
    Builtin(NativeFn),
    Bool(bool),
}

impl Display for Value {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Value::Atom(a) => write!(f, "Atom::{a}"),
            
            Value::Str(s) => write!(f, "{s}"),
            
            Value::Number(i) => write!(f, "{i}"),
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