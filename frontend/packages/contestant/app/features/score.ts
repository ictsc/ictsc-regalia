import type { Score as ProtoScore } from "@ictsc/proto/contestant/v1";
import type { ComponentProps } from "react";
import type { Score } from "../components/score";

export function protoScoreToProps(
  maxScore: number,
  protoScore?: ProtoScore,
): ComponentProps<typeof Score> {
  return {
    maxScore: maxScore,
    score: protoScore?.score,
    rawScore: protoScore?.markedScore,
    penalty: protoScore?.penalty,
    fullScore: protoScore?.score === maxScore,
    rawFullScore: protoScore?.markedScore === maxScore,
  };
}
