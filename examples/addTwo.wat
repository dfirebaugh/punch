(module
    (func $addTwo (export "addTwo")(param $x i32)(param $y i32)(result i32)
        (return (i32.add (local.get $x) (local.get $y)))
    )
)
