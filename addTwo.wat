(module
    (func $addTwo (export "addTwo")(param $x i32)(param $y i32)(result i32)
        (return (i32.add (local.get $x) (local.get $y)))
    )

    (func $addFour (export "addFour")(param $a i32)(param $b i32)(param $c i32)(param $d i32)(result i32)
        (return (i32.add (i32.add (i32.add (local.get $a) (local.get $b)) (local.get $c)) (local.get $d)))
    )
)
