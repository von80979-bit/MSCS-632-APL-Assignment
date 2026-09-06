// Rust rejects this file. The compiler message is the result the program exists to show, so there is no binary to
// run and nothing to profile.
fn main() {
    let original_values = vec![1, 2, 3, 4, 5];
    let moved_values = original_values;  // ownership moves here, so original_values stops being usable
    println!("sum after the move: {}", moved_values.iter().sum::<i32>());
    println!("the moved-from binding: {:?}", original_values);
}
