import os
import sys
from pathlib import Path

import orjson
import ray


def error(text: str) -> None:
    print(f"error: {text}", file=sys.stderr)
    help()
    sys.exit(1)


def help() -> None:
    print("usage:\n main.py logs.json 4")


@ray.remote(max_retries=3)
def analyze(filename: Path | str, start: int, end: int) -> tuple[int, int]:
    """
    analyze parses JSON log-rows in filename and calculates a running sum
    of the `temperature` field.

    It works with the following format:
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

    total = 0
    count = 0
    with open(filename, "rb") as rfile:
        if start != 0:
            # skip a partial line (the previous chunk will handle it)
            rfile.seek(start - 1)
            rfile.readline()

        pos = rfile.tell()
        while pos < end:
            line = rfile.readline()
            if not line:
                break

            # will absolutely shit itself on invalid JSON data
            dump = orjson.loads(line)
            total += dump["temperature"]
            count += 1

            pos += len(line)

    return (total, count)


DEFAULT_NUM_CPUS: int = os.cpu_count() or 1


def main(filename: Path, num_cpus=DEFAULT_NUM_CPUS) -> None:
    print(f"starting ray with {num_cpus=}")
    ray.init(num_cpus=num_cpus)

    chunks = num_cpus * 20  # tune this knob, but 20 is more or less fine

    size = filename.stat().st_size
    bounds = [size * i // chunks for i in range(chunks + 1)]

    refs = [analyze.remote(filename, bounds[i], bounds[i + 1]) for i in range(chunks)]
    results = ray.get(refs)

    total = sum(t for t, _ in results)
    count = sum(c for _, c in results)

    print(f"the average temperature is {(total / count):.2f}")


ARGUMENT_COUNT: int = 2

if __name__ == "__main__":
    if len(sys.argv[1:]) != ARGUMENT_COUNT:
        error(f"requires {ARGUMENT_COUNT} arguments, got {len(sys.argv[1:])}")

    filename = Path(sys.argv[1])
    if not filename.is_file():
        error(f"{filename} is not a regular file")

    try:
        num_cpus = int(sys.argv[2])

    except ValueError:
        error(f"{sys.argv[2]} is not an int")

    main(filename, num_cpus)  # type: ignore
