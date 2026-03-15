"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.protoScoreToProps = protoScoreToProps;
function protoScoreToProps(maxScore, protoScore) {
    return {
        maxScore: maxScore,
        score: protoScore === null || protoScore === void 0 ? void 0 : protoScore.score,
        rawScore: protoScore === null || protoScore === void 0 ? void 0 : protoScore.markedScore,
        penalty: protoScore === null || protoScore === void 0 ? void 0 : protoScore.penalty,
        fullScore: (protoScore === null || protoScore === void 0 ? void 0 : protoScore.score) === maxScore,
        rawFullScore: (protoScore === null || protoScore === void 0 ? void 0 : protoScore.markedScore) === maxScore,
    };
}
