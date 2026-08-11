use log::{trace, debug};

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
                Value::Lit(Literal::Num(i)) => i,
                v => {
                    return Err(format!("Expected string, got {v:?}"))
                }
            };
        let lhs =
        match e.stack.pop().unwrap() {
            Value::Int(i) => i,
            Value::Lit(Literal::Num(i)) => i,
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
        if let Some(Value::Bool(_)) = e.stack.pop() {
            e.stack.push(Value::Bool(true))
        } else {
            e.stack.push(Value::Bool(false))
        }

        Ok(())
    }
}
builtin! {
    is_number, e, {
        if let Some(Value::Int(_)) = e.stack.pop() {
            e.stack.push(Value::Bool(true))
        } else if let Some(Value::Float(_)) = e.stack.pop() {
            e.stack.push(Value::Bool(true))
        } else {
            e.stack.push(Value::Bool(false))
        }

        Ok(())
    }
}
builtin! {
    is_int, e, {
        if let Some(Value::Int(_)) = e.stack.pop() {
            e.stack.push(Value::Bool(true))
        } else {
            e.stack.push(Value::Bool(false))
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