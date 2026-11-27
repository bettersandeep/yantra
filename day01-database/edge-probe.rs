use yantra_sql::Database;
#[test]
fn probe() {
    let qs = [
        "SELECT (-9223372036854775807 - 1) % -1",
        "SELECT (-9223372036854775807 - 1) / -1",
        "SELECT - (-9223372036854775807 - 1)",
        "SELECT 7 / -2, -7 / 2, 7 % -2, -7 % 2",
        "select 2*3 col0, 4 AS col1",
        "SELECT + + - - 5",
        "SELECT 4 / 0 + 1, 0 * (1/0), 1 % 0",
        "SELECT 00012",
        "SELECT 9223372036854775807 + 0",
        "SELECT 9223372036854775807 + 1",
        "SELECT 1 /* c */ + 1",
        "SELECT DISTINCT ALL 1",
        "SELECT ALL DISTINCT 1",
        "SELECT 1 +", "SELECT", "SELECT 1,", "SELECT (1", "SELECT 1)", "SELECT 1 FROM", "SELECT 2 2",
        "SELECT 1 AS", "SELECT 1 AS 2", "SELECT 1 AS select", "SELECT 1e3", "SELECT 1.", "SELECT .5",
        "SELECT\u{a0}1", "SELECT 1 x y", "SELECT 1 AS x, 2 AS x", "", "   ", "SELECT 1;", "SELECT 1 AS \"x\"",
        "SELECT 1 AS x1_y", "SELECT 1 AS _x", "SELECT 1 AS 1x", "SELECT 1x", "SELECT 0x10", "SELECT 1 - -1", "SELECT 1--1",
        "SELECT 2 * 3 % 4 / 5 - 6 + 7", "SELECT 2 - 3 - 4", "SELECT 100 / 7 / 2", "SELECT -(2+3)*-4",
    ];
    for q in qs {
        let r = std::panic::catch_unwind(|| Database::new().execute(q));
        let s = match r { Err(_) => "PANIC".to_string(), Ok(Err(e)) => format!("error: {e}"), Ok(Ok(rows)) => format!("{rows:?}") };
        println!("{q:45} => {s}");
    }
}
