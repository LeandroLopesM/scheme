use std::io::{Write, stdout};

use log::{debug, trace, warn};

use crate::scheme::{engine::{Engine, Value}, parser::Literal};

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
                Value::Int(i) => i,
                Value::Lit(Literal::Int(i)) => i,
                v => {
                    return Err(format!("Expected string, got {v:?}"))
                }
            };
        let lhs =
        match e.stack.pop().unwrap() {
            Value::Int(i) => i,
            Value::Lit(Literal::Int(i)) => i,
            v => {
                return Err(format!("Expected string, got {v:?}"))
            }
        };

        trace!("Result was: {}", lhs + rhs);

        e.stack.push(Value::Int(lhs + rhs));

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