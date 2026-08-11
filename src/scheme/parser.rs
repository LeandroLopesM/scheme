use std::{fmt::Display, iter::Peekable};

use log::trace;

use crate::scheme::parser::TokenKind::Ident;

#[derive(Debug)]
pub struct Error {
    msg: String,
    line: u32,
}

impl Display for Error {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}: {}", self.line, self.msg)
    }
}

impl Error {
    pub fn new<S: ToString>(line: u32, msg: S) -> Self {
        Self {
            line,
            msg: msg.to_string(),
        }
    }
}

#[derive(Clone, Debug, PartialEq)]
#[allow(dead_code)]
pub enum Literal {
    Num(i64),
    Float(f64),
    Str(String),
    Bool(bool),
}

#[derive(Clone, Debug)]
#[allow(dead_code)]
pub enum TokenKind {
    Ident(String),
    Literal(Literal),
    Scope(Scope),
}

#[derive(Clone, Debug)]
#[allow(dead_code)]
pub struct Token {
    pub line: u32,
    pub kind: TokenKind,
}

#[derive(Clone, Debug)]
#[allow(dead_code)]
pub struct Scope {
    pub line: u32,
    pub name: String,
    pub args: Vec<Token>,
}

impl Token {
    pub fn new(kind: TokenKind, line: u32) -> Self {
        Self { line, kind }
    }
}

pub fn lex(str: String) -> Result<Vec<Scope>, Error> {
    let mut out = Vec::new();
    let mut iter = str.chars().peekable();
    let mut line = 1;

    while let Some(c) = iter.next() {
        match c {
            ';' => {
                parse_comm(&mut iter);
                line += 1;
            }
            '(' => match parse_scope(&mut iter, &mut line) {
                Ok(s) => out.push(s),
                Err(e) => return Err(e),
            },

            ' ' => continue,
            '\n' => line += 1,
            _ => {
                return Err(Error::new(
                    line,
                    format!("Invalid character in global scope: {c}"),
                ));
            }
        }
    }

    Ok(out)
}

pub fn parse_scope(
    mut iter: &mut Peekable<impl Iterator<Item = char>>,
    line: &mut u32,
) -> Result<Scope, Error> {
    let mut out: Vec<Token> = Vec::new();
    let start = *line;

    while let Some(' ') = iter.peek() {
        _ = iter.next()
    }
    let name = if let Ident(i) = parse_id(iter) {
        i
    } else {
        unreachable!()
    };

    while let Some(c) = iter.peek() {
        match c {
            '\n' => *line += 1,

            '"' => {
                _ = iter.next(); // purge '"'

                match parse_str(iter, *line) {
                    Ok(str) => out.push(Token::new(str, *line)),
                    Err(e) => return Err(e),
                }
            }

            ';' => {
                parse_comm(&mut iter);
                *line += 1;
            }

            'a'..'z' | 'A'..'Z' | '_' => {
                out.push(Token::new(parse_id(&mut iter), *line));
            }

            '#' => {
                _  = iter.next();

                if let Some('t') = iter.next() {
                    out.push(Token::new(TokenKind::Literal(Literal::Bool(true)), *line))
                } else if let Some('f') = iter.next() {
                    out.push(Token::new(TokenKind::Literal(Literal::Bool(false)), *line))
                } else {
                    return Err(Error::new(*line, "Invalid boolean literal"));
                }
            }

            '0'..'9' => {
                match parse_num(&mut iter, *line) {
                    Ok(i) => out.push(Token::new(i, *line)),
                    Err(e) => return Err(e),
                };
            }

            '(' => {
                _ = iter.next(); // skip '('

                match parse_scope(iter, line) {
                    Ok(s) => out.push(Token::new(TokenKind::Scope(s), *line)),
                    Err(e) => return Err(e),
                }
            }

            ')' => {
                _ = iter.next(); // purge ')'

                return Ok(Scope {
                    line: start,
                    name,
                    args: out,
                });
            }

            c => {
                trace!("Ignored '{c}'");
                _ = iter.next();
            }
        }
    }

    // if let None = iter.peek() {
    Err(Error {
        line: start,
        msg: format!("Unclosed scope {name}"),
    })
    // } else {
    //     Ok(Scope {
    //         line: *line,
    //         name: name,
    //         args: out,
    //     })
    // }
}

pub fn parse_comm(iter: &mut Peekable<impl Iterator<Item = char>>) {
    while let Some(c) = iter.next()
        && c != '\n'
    {}
}
pub fn parse_id(iter: &mut Peekable<impl Iterator<Item = char>>) -> TokenKind {
    let mut out = String::new();

    while let Some(&c) = iter.peek() {
        if c.is_whitespace() || c == '(' || c == ')' {
            break;
        }

        out.push(c);
        _ = iter.next();
    }

    trace!("Ident {out}");
    TokenKind::Ident(out)
}

fn parse_num(
    iter: &mut Peekable<impl Iterator<Item = char>>,
    line: u32,
) -> Result<TokenKind, Error> {
    let mut n = Literal::Num(0);
    let mut buffer = String::new();
    let mut radix = 10;

    if buffer.starts_with("0x") {
        radix = 16;
        buffer = String::from(&buffer[2..])
    } else if buffer.starts_with("0b") {
        radix = 2;
        buffer = String::from(&buffer[2..])
    }

    if buffer.contains(".") {
        n = Literal::Float(0.)
    }

    while let Some(&c) = iter.peek() {
        if !c.is_numeric() && !('A'..'G').contains(&c) && !('a'..'G').contains(&c) {
            break;
        }

        buffer.push(c);
        _ = iter.next();
    }

    if let Literal::Num(_) = n {
        trace!("Intlit: ({radix}) {buffer}");

        match i64::from_str_radix(&buffer, radix) {
            Ok(n) => return Ok(TokenKind::Literal(Literal::Num(n))),
            Err(e) => {
                return Err(Error {
                    msg: format!("Invalid {} number '{}'\n{e}", readable_radix(radix), buffer),
                    line: line,
                });
            }
        };
    } else {
        if radix != 10 {
            return Err(Error {
                msg: format!(
                    "Number cannot be float and {} ({buffer})",
                    readable_radix(radix)
                ),
                line: line,
            });
        }

        trace!("Floatlit: {buffer}");

        match buffer.parse::<f64>() {
            Ok(n) => return Ok(TokenKind::Literal(Literal::Float(n))),
            Err(e) => {
                return Err(Error {
                    msg: format!("Invalid {} number '{}'\n{e}", readable_radix(radix), buffer),
                    line: line,
                });
            }
        };
    }
}

fn parse_str(
    iter: &mut Peekable<impl Iterator<Item = char>>,
    line: u32,
) -> Result<TokenKind, Error> {
    let mut buffer = String::new();
    let mut done = false;

    while let Some(&c) = iter.peek() {
        if c == '"' {
            done = true;
            break;
        }

        buffer.push(c);
        _ = iter.next()
    }

    trace!("StringLit: \"{buffer}\"");

    if !done {
        Err(Error::new(line, "Unclosed string"))
    } else {
        _ = iter.next(); // purge '"'
        Ok(TokenKind::Literal(Literal::Str(buffer)))
    }
}

fn readable_radix(radix: u32) -> &'static str {
    match radix {
        16 => "hexadecimal",
        2 => "binary",
        10 => "decimal",
        _ => unreachable!(),
    }
}
