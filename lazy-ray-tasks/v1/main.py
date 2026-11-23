import sys
from pathlib import Path

import orjson


def error(text: str) -> None:
    print(text)
    help()
    sys.exit(1)


def help() -> None:
    print("usage:\n main.py logs.json")


class LogAnalyzer:
    """
    LogAnalyzer works with JSON logs in the following format:
    ```
    {
        "time": "2026-11-17T13:33:56.7246204+02:00",
        "level": "INFO",
        "msg": "recorded temperature",
        "temperature": 62, // int [1; 100]
        "instance_id": "5f800d4d-dd55-4624-adc5-68ad54201211" // UUIDv4
    }
    ```
    """

    count: int
    sum: int

    def __init__(self) -> None:
        self.count = 0
        self.sum = 0

    def analyze(self, filename: Path | str, start: int, end: int) -> float:
        """
        analyze parses JSON log-rows in filename and calculates a running sum
        of the `temperature` field.
        """

        with open(filename, "rb") as rfile:
            for raw in rfile:
                # will absolutely shit itself on non-JSON data
                dump = orjson.loads(raw)
                self.sum += dump["temperature"]
                self.count += 1

        return self.sum / self.count


def main(filename: Path | str) -> None:
    l = LogAnalyzer()
    l.analyze(filename, 0, 0)

    print(f"count={l.count}\nsum={l.sum}\naverage={(l.sum / l.count):.4f}")


ARGUMENT_COUNT = 1

if __name__ == "__main__":
    if len(sys.argv[1:]) != ARGUMENT_COUNT:
        error(f"error: requires {ARGUMENT_COUNT} arguments, got {len(sys.argv[1:])}")

    filename = Path(sys.argv[1])
    main(filename)
