use std::io::{Write, stdout};

use log::{debug, trace, warn};

use crate::scheme::{datatypes::{Number, Value}, engine::Engine, parser::Literal};

macro_rules! builtin {
    ($name:ident, $eng:ident, $ct:block) => {
        pub fn $name($eng: &mut Engine) -> Result<(), String> {
            debug!("BICall '{}'\nStack: {:?}", stringify!($name), $eng.stack);
            $ct
        }
    };
}

builtin! {
    display, e, {
        let val = e.stack.pop().unwrap();

        print!("{val}");

        if e.config.display_flushes {
            if let Err(err) = stdout().flush() {
                warn!("Engine configuration failed, unsetting CFG::DisplayFlushes\nError: {err}");
                e.config.display_flushes = false;
            }
        }

        Ok(())
    }
}

builtin!{
    newline, _e, {
        Ok(println!())
    }
}


builtin!{
    add, e, {
        let rhs =
            match e.stack.pop().unwrap() {
                Value::Number(n) => n,
                Value::Lit(Literal::Int(i)) => Number::Int(i),
                v => {
                    return Err(format!("Expected number, got {v:?}"))
                }
            };
        let lhs =
        match e.stack.pop().unwrap() {
            Value::Number(i) => i,
            Value::Lit(Literal::Int(i)) => Number::Int(i),
            v => {
                return Err(format!("Expected number, got {v:?}"))
            }
        };

        e.stack.push(Value::Number(lhs + rhs));

        Ok(())
    }
} 

builtin! {
    is_bool, e, {
        match e.stack.pop() {
            Some(Value::Bool(_)) => e.stack.push(Value::Bool(true)),
            Some(Value::Lit(Literal::Bool(_))) => e.stack.push(Value::Bool(true)),
            _ => e.stack.push(Value::Bool(false))
        }

        Ok(())
    }
}
builtin! {
    is_number, e, {
        match e.stack.pop() {
            Some(Value::Int(_)) => e.stack.push(Value::Bool(true)),
            Some(Value::Lit(Literal::Int(_))) => e.stack.push(Value::Bool(true)),
            
            Some(Value::Float(_)) => e.stack.push(Value::Bool(true)),
            Some(Value::Lit(Literal::Float(_))) => e.stack.push(Value::Bool(true)),

            _ => e.stack.push(Value::Bool(false))
        }

        Ok(())
    }
}
builtin! {
    is_int, e, {
        match e.stack.pop() {
            Some(Value::Int(_)) => e.stack.push(Value::Bool(true)),
            Some(Value::Lit(Literal::Int(_))) => e.stack.push(Value::Bool(true)),

            _ => e.stack.push(Value::Bool(false))
        }

        Ok(())
    }
}
builtin! {
    eqv, e, {
        let rhs = e.stack.pop().unwrap();
        let lhs = e.stack.pop().unwrap();

        e.stack.push(Value::Bool(rhs == lhs));

        Ok(())
    }
}