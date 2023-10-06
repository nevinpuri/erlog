/// gramar
/// field.hi = 'fd' or field.bye
/// field.hi > 3.0
/// 'asdf' in field.hi
/// has(field.hi)
use nom::{
    branch::alt,
    bytes::complete::{tag, take_till, take_until, take_while, take_while1},
    character::{is_digit, is_space},
    combinator::success,
    error::Error,
    IResult,
};
use pyo3::{
    exceptions::PyValueError, pyclass, pyfunction, pymodule, types::PyModule, wrap_pyfunction,
    FromPyObject, IntoPy, Py, PyAny, PyResult, Python,
};

#[derive(Debug, Clone)]
#[pyclass]
pub enum Operation {
    Eq,
    Lt,
    Lte,
    Gt,
    Gte,
}

#[derive(Debug)]
#[pyclass]
pub struct ParseOut {
    #[pyo3(get)]
    pub and: Vec<Expression>,

    #[pyo3(get)]
    pub or: Vec<Expression>,

    #[pyo3(get)]
    pub none: Option<Expression>,
}

#[derive(Debug, Clone)]
#[pyclass]
pub struct Expression {
    #[pyo3(get)]
    pub field: String,
    #[pyo3(get)]
    pub operation: Operation,
    #[pyo3(get)]
    pub val: Val,
}

#[derive(Debug, FromPyObject, Clone)]
pub enum Val {
    Str(String),
    Num(f32),
    Bool(bool),
}

impl IntoPy<Py<PyAny>> for Val {
    fn into_py(self, py: Python<'_>) -> Py<PyAny> {
        match self {
            Self::Bool(e) => e.into_py(py),
            Self::Num(e) => e.into_py(py),
            Self::Str(e) => e.into_py(py),
        }
    }
}

fn get_char_from_op(op: &Operation) -> String {
    match op {
        Operation::Eq => "=",
        Operation::Lt => "<",
        Operation::Lte => "<=",
        Operation::Gt => ">",
        Operation::Gte => ">=",
    }
    .into()
}

fn get_eq(input: &str) -> IResult<&str, (&str, Operation)> {
    let (input, field) = take_until("=")(input)?;
    Ok((input, (field, Operation::Eq)))
}

fn get_lt(input: &str) -> IResult<&str, (&str, Operation)> {
    let (input, field) = take_until("<")(input)?;
    Ok((input, (field, Operation::Lt)))
}

fn get_lte(input: &str) -> IResult<&str, (&str, Operation)> {
    let (input, field) = take_until("<=")(input)?;
    Ok((input, (field, Operation::Lte)))
}

fn get_gt(input: &str) -> IResult<&str, (&str, Operation)> {
    let (input, field) = take_until(">")(input)?;
    Ok((input, (field, Operation::Gt)))
}

fn get_gte(input: &str) -> IResult<&str, (&str, Operation)> {
    let (input, field) = take_until(">=")(input)?;
    Ok((input, (field, Operation::Gte)))
}
// fn get_and(input: &str) -> IResult<

fn parse_string(input: &str) -> IResult<&str, Val> {
    let (input, _) = tag("\'")(input)?;
    let (input, val) = take_until("\'")(input)?;
    let (input, _) = tag("\'")(input)?;

    Ok((input, Val::Str(val.into())))
}

fn parse_num(input: &str) -> IResult<&str, Val> {
    let (input, val) = take_till(|c: char| c.is_whitespace())(input)?;
    let val: f32 = match val.parse::<f32>() {
        Ok(val) => val,
        Err(e) => {
            return Err(nom::Err::Error(Error::new(
                input,
                nom::error::ErrorKind::Fail,
            )));
        }
    };

    Ok((input, Val::Num(val)))
}

fn parse_bool(input: &str) -> IResult<&str, Val> {
    let (input, val) = take_till(|c: char| c.is_whitespace())(input)?;
    let converted = match val {
        "true" => true,
        "false" => false,
        _ => {
            return Err(nom::Err::Error(Error::new(
                input,
                nom::error::ErrorKind::Fail,
            )));
        }
    };

    Ok((input, Val::Bool(converted)))
}

/// parses the corresponding value after the field
fn parse_val(input: &str) -> IResult<&str, Val> {
    let (input, _) = take_while1(char::is_whitespace)(input)?;
    let (input, value) = take_till(char::is_whitespace)(input)?;
    let (_, val) = alt((parse_string, parse_num, parse_bool))(value)?;
    Ok((input, val))
}

/// gets the field and operation from input
fn get_field_op(input: &str) -> IResult<&str, (&str, Operation)> {
    let (input, (field, op)) = alt((get_lte, get_gte, get_lt, get_gt, get_eq))(input)?;
    let c = get_char_from_op(&op);
    let field = field.trim();
    let (input, _) = tag(c.as_str())(input)?;

    Ok((input, (field, op)))
}

fn parse_expression(input: &str) -> IResult<&str, Expression> {
    let (input, (field, op)) = get_field_op(input)?;
    let (input, val) = parse_val(input)?;

    Ok((
        input,
        Expression {
            field: field.into(),
            operation: op,
            val,
        },
    ))
}

#[derive(Debug, Clone, PartialEq, PartialOrd)]
enum Conditional {
    And,
    Or,
    None,
}

fn match_and(input: &str) -> IResult<&str, Conditional> {
    let (input, _) = alt((tag("AND"), tag("and")))(input)?;
    Ok((input, Conditional::And))
}

fn match_or(input: &str) -> IResult<&str, Conditional> {
    let (input, _) = alt((tag("OR"), tag("or")))(input)?;
    Ok((input, Conditional::Or))
}

fn match_none(input: &str) -> IResult<&str, Conditional> {
    let (input, c) = success(Conditional::None)(input)?;
    Ok((input, c))
}

fn parse_conditional(input: &str) -> IResult<&str, Conditional> {
    let (input, _) = take_while(|c: char| c.is_whitespace())(input)?;
    let (input, op) = alt((match_and, match_or, match_none))(input)?;

    Ok((input, op))
}

fn parse_until_next(input: &str) -> IResult<&str, (Expression, Conditional)> {
    println!("parsing until next");
    let (input, expr) = parse_expression(input)?;
    let (input, c) = parse_conditional(input)?;

    Ok((input, (expr, c)))
}

fn _parse_input(input: &str) -> IResult<&str, ParseOut> {
    // let mut ands = vec![];
    let mut p = ParseOut {
        and: vec![],
        or: vec![],
        none: None,
    };

    let mut a = parse_until_next(input)?;
    // let (mut input, (mut expr, mut c)) = parse_until_next(input)?;

    if a.1 .1 == Conditional::None {
        p.none = Some(a.1 .0);
        return Ok((a.0, p));
    }

    let mut prev_op = a.1 .1.clone();

    let f = loop {
        let expr = a.1 .0.clone();

        match a.1 .1.clone() {
            Conditional::And => p.and.push(expr),
            Conditional::Or => p.or.push(expr),
            _ => unreachable!(),
            // Conditional::None => {
            //     break a.clone();
            // }
        }
        // p.and.push(expr);

        prev_op = a.1 .1.clone();

        a = parse_until_next(a.0)?;

        if a.1 .1 == Conditional::None {
            break a.clone();
        }
    };

    match prev_op {
        Conditional::And => p.and.push(f.1 .0),
        Conditional::Or => p.or.push(f.1 .0),
        _ => p.none = Some(f.1 .0),
    }

    // println!("{:#?} - {:#?}", f, prev_op);

    // figure out what previous shit was and then add to that array
    // if it is none, then just push to fianl expr or something

    // let (input, expr) = parse_expression(input)?;
    // let (input, c) = parse_conditional(input)?;

    // Ok((
    //     input,
    //     ParseOut {
    //         and: ands,
    //         or: vec![],
    //         none: Expression {
    //             field: "".into(),
    //             operation: Operation::Eq,
    //             val: "".into(),
    //         },
    //     },
    // ))

    Ok((a.0, p))
    // let (input, (field, op)) = get_field_op(input)?;
    // let (input, val) = parse_val(input)?;

    // Ok((
    //     input,
    //     Expression {
    //         field: field.into(),
    //         operation: op,
    //         val,
    //     },
    // ))
}

#[pyfunction]
fn parse_input(input: &str) -> PyResult<ParseOut> {
    if input == "" {
        return Ok(ParseOut {
            and: vec![],
            or: vec![],
            none: None,
        });
    }

    let (input, out) = match _parse_input(input) {
        Ok(e) => e,
        Err(e) => return Err(PyValueError::new_err("error parsing user input")),
    };

    println!("rust out {:#?}", out);
    return Ok(out);
}

#[pymodule]
fn erlog_utils(_py: Python, m: &PyModule) -> PyResult<()> {
    m.add_function(wrap_pyfunction!(parse_input, m)?)?;
    m.add_class::<Expression>()?;
    m.add_class::<ParseOut>()?;
    m.add_class::<Operation>()?;

    Ok(())
}

#[cfg(test)]
mod test {
    use super::*;
    #[test]
    fn test_main() {
        let out = parse_input("field.hi = 'fadsfadf'").unwrap();
        println!("{:#?}", out);
        println!("Hello, world!");
    }
}

// fn main() {
//     let (input, out) = parse_input("field.hi = 'fadsfadf'").unwrap();
//     println!("{:#?}", out);
//     println!("Hello, world!");
// }
