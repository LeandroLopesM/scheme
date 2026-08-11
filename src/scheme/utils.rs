use crate::scheme::parser::{Scope, Token};

use super::parser::TokenKind;

pub fn describe(t: Vec<Scope>, lvl: usize) -> String {
    let mut o = String::new();

    for e in t {
        o += &format!("\n{}Scope {}", "| ".repeat(lvl), e.name);
        o += &describe_tok(e.args, lvl + 1);
    }

    o
}

pub fn describe_tok(toks: Vec<Token>, lvl: usize) -> String {
    let mut o = String::new();

    for tok in toks {
        match tok.kind {
            TokenKind::Scope(s) => o += &describe(vec![s], lvl + 1),
            f => o += &format!("\n{}{:?}", "| ".repeat(lvl), f),
        }
    }

    o
}
