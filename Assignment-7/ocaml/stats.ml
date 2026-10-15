(* Statistics calculator in functional OCaml: immutable int lists and higher-order functions, with no mutable state.
   Build: ocamlopt -o stats stats.ml *)

let count values = List.length values

let sum values = List.fold_left ( + ) 0 values

let mean values = float_of_int (sum values) /. float_of_int (count values)

let median values =
  let sorted = List.sort compare values in
  let n = count values in
  if n mod 2 = 1 then float_of_int (List.nth sorted (n / 2))
  else float_of_int (List.nth sorted ((n / 2) - 1) + List.nth sorted (n / 2)) /. 2.0

(* Equal values sit next to each other in the sorted list, so one fold builds the (value, occurrences) pairs. *)
let occurrences values =
  let add_value pairs value =
    match pairs with
    | (last, n) :: rest when last = value -> (last, n + 1) :: rest
    | _ -> (value, 1) :: pairs
  in
  List.rev (List.fold_left add_value [] (List.sort compare values))

let mode values =
  let pairs = occurrences values in
  let highest = List.fold_left (fun best (_, n) -> max best n) 0 pairs in
  pairs |> List.filter (fun (_, n) -> n = highest) |> List.map fst

let min_value values = List.fold_left min (List.hd values) values

let max_value values = List.fold_left max (List.hd values) values

let range values = max_value values - min_value values

let format_list values = "[" ^ String.concat ", " (List.map string_of_int values) ^ "]"

let print_report list_number values =
  Printf.printf "List %d: %s\n" list_number (format_list values);
  Printf.printf "  Count:  %d\n" (count values);
  Printf.printf "  Mean:   %.2f\n" (mean values);
  Printf.printf "  Median: %.2f\n" (median values);
  Printf.printf "  Mode:   %s\n" (format_list (mode values));
  Printf.printf "  Min:    %d\n" (min_value values);
  Printf.printf "  Max:    %d\n" (max_value values);
  Printf.printf "  Range:  %d\n" (range values)

let demo_lists = [ [ 4; 8; 15; 16; 23; 42; 8 ]; [ 7; 3; 9; 3; 7; 1; 12; 2 ]; [ 5; -2; 10; 0; -7 ] ]

let () =
  print_endline "=== Statistics Calculator: OCaml (functional) ===";
  List.iteri
    (fun index values ->
      print_newline ();
      print_report (index + 1) values)
    demo_lists
