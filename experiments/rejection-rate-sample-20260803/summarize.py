#!/usr/bin/env python3
"""Deterministically derive the reported rejection and ingest-cost figures."""

from __future__ import annotations

import argparse
import json
from decimal import Decimal
from pathlib import Path


HERE = Path(__file__).resolve().parent


def decimal(value: object) -> Decimal:
    return Decimal(str(value))


def summarize(source: Path) -> dict[str, object]:
    data = json.loads(source.read_text(encoding="utf-8"))
    observation = data["observation"]
    pricing = data["pricing_assumption"]
    tallies = data["tallies"]

    documents = decimal(observation["documents"])
    candidates = decimal(tallies["candidates"])
    prompt_tokens = decimal(observation["prompt_tokens"])
    completion_tokens = decimal(observation["completion_tokens"])
    gate_rejections = sum(
        decimal(tallies[key])
        for key in (
            "rejected_empty_or_overlong_quote",
            "rejected_bad_target_range",
            "rejected_quote_not_in_target",
            "rejected_generic_speaker",
        )
    )
    deduplicated = decimal(tallies["deduplicated"])
    sample_cost = (
        prompt_tokens * decimal(pricing["prompt_usd_per_million_tokens"])
        + completion_tokens
        * decimal(pricing["completion_usd_per_million_tokens"])
    ) / Decimal(1_000_000)

    return {
        "sample_documents": int(documents),
        "candidate_claims": int(candidates),
        "gate_rejections": int(gate_rejections),
        "gate_rejection_rate_percent": f"{gate_rejections / candidates * 100:.6f}",
        "gate_rejection_rate_percent_reported": f"{gate_rejections / candidates * 100:.2f}",
        "overall_drops_including_deduplication": int(gate_rejections + deduplicated),
        "overall_drop_rate_percent": f"{(gate_rejections + deduplicated) / candidates * 100:.6f}",
        "validated_claims_per_document": f"{decimal(tallies['validated']) / documents:.2f}",
        "prompt_tokens_per_document": f"{prompt_tokens / documents:.2f}",
        "completion_tokens_per_document": f"{completion_tokens / documents:.2f}",
        "total_tokens_per_document": f"{(prompt_tokens + completion_tokens) / documents:.2f}",
        "sample_cost_usd": f"{sample_cost:.7f}",
        "cost_usd_per_document": f"{sample_cost / documents:.8f}",
        "cost_usd_per_document_reported": f"{sample_cost / documents:.3f}",
        "five_hundred_document_token_extrapolation": int(
            (prompt_tokens + completion_tokens) / documents * 500
        ),
        "five_hundred_document_cost_usd_extrapolation": f"{sample_cost / documents * 500:.5f}",
    }


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--input",
        type=Path,
        default=HERE / "observed-result.json",
        help="captured replay observation",
    )
    parser.add_argument(
        "--check",
        type=Path,
        help="fail unless the derived JSON exactly matches this file",
    )
    args = parser.parse_args()

    result = summarize(args.input)
    rendered = json.dumps(result, indent=2, sort_keys=True) + "\n"
    if args.check:
        expected = args.check.read_text(encoding="utf-8")
        if rendered != expected:
            raise SystemExit(f"derived output differs from {args.check}")
    print(rendered, end="")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
