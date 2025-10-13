import io
import unittest
from main import solve, run


def trim_lines(text: str) -> str:
    lines = text.splitlines()
    # Убираем trailing whitespace с каждой строки
    lines = [line.rstrip() for line in lines]
    # Убираем пустые строки в конце
    while lines and lines[-1] == "":
        lines.pop()
    return "\n".join(lines)


class TestRunSolve(unittest.TestCase):

    def test_run_solve(self):
        self._test_run(solve)

    def _test_run(self, solve_func):
        test_cases = [
            # (name, input_str, expected_output_str)
            ("1", "3 3\n1 2 3\n6 5 4\n7 8 9\n", "9"),
            ("2", "3 3\n2 2 2\n2 3 2\n2 2 2\n", "2"),
            ("3", "1 1\n1", "1"),
            ("4", "3 3\n2 4 5\n2 3 4\n2 2 2\n", "4"),
            ("5_all_identical", "3 4\n5 5 5 5\n5 5 5 5\n5 5 5 5\n", "1"),
            ("6_full_sequence_1_16", "4 4\n1 2 3 4\n8 7 6 5\n9 10 11 12\n16 15 14 13\n", "16"),
            ("7_chain_length_2", "2 2\n1 3\n2 4\n", "2"),
            ("8_long_vertical_chain", "5 2\n1 2\n3 4\n5 6\n7 8\n9 10\n", "2"),
            ("9_long_horizontal_chain", "2 5\n1 2 3 4 5\n10 9 8 7 6\n", "10"),
            ("10_chain_with_zero", "3 3\n0 1 2\n5 4 3\n6 7 8\n", "9"),
            ("11_multiple_chains", "3 3\n1 2 1\n4 3 4\n5 6 5\n", "6"),
            ("12_single_row_sequence", "1 5\n1 2 3 4 5\n", "5"),
            ("13_single_column_sequence", "5 1\n1\n2\n3\n4\n5\n", "5"),
            ("14_zigzag_chain", "3 3\n1 3 2\n4 6 5\n7 9 8\n", "2"),
            ("15_max_chain_not_from_max_value", "3 3\n9 1 2\n8 7 3\n5 4 6\n", "3"),
            ("16_large_numbers_long_chain", "2 3\n1000000 1000001 1000002\n1000005 1000004 1000003\n", "6"),
            ("17_isolated_chains", "3 3\n1 5 2\n4 3 6\n7 8 9\n", "3"),
            ("18_obstacle_in_middle", "3 4\n1 2 3 4\n10 9 8 5\n11 12 7 6\n", "12"),
            ("19_all_decreasing", "3 3\n9 8 7\n6 5 4\n3 2 1\n", "3"),
            ("20_single_chain_length_2", "3 3\n1 1 1\n1 2 1\n1 1 1\n", "2"),
            ("21_corner_chain", "3 3\n1 2 1\n1 1 1\n1 1 1\n", "2"),
            ("22_two_equal_chains", "3 3\n1 2 4\n3 5 5\n4 3 2\n", "3"),
            ("23_spiral_chain", "3 3\n1 2 3\n8 9 4\n7 6 5\n", "9"),
            ("24_minimal_2x2_with_chain", "2 2\n1 2\n4 3\n", "4"),
            ("25_single_chain_continuation", "3 3\n1 1 1\n1 2 1\n1 3 1\n", "3"),
        ]

        for name, input_str, want_out in test_cases:
            with self.subTest(name=name):
                in_stream = io.StringIO(input_str)
                out_stream = io.StringIO()

                run(in_stream, out_stream)

                got_out = out_stream.getvalue()
                self.assertEqual(trim_lines(got_out), trim_lines(want_out))


if __name__ == "__main__":
    unittest.main()
