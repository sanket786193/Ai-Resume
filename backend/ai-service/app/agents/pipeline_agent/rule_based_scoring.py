"""Rule-based ATS scoring (no LLM)."""
from typing import Any, Dict


class RuleBasedScoringEngine:
    def calculate_skills_score(self, required: list, candidate: list) -> float:
        if not required:
            return 0.0
        matched = len(set(required).intersection(set(candidate)))
        return min((matched / len(required)) * 30, 30)

    def calculate_experience_score(self, required_years: float, candidate_years: float) -> float:
        if required_years <= 0:
            return 0.0
        return min((candidate_years / required_years) * 25, 25)

    def calculate_ats_score(self, ats: Dict[str, Any]) -> float:
        score = 0.0
        if ats.get("keyword_density_ok"):
            score += 5
        if ats.get("standard_headings"):
            score += 4
        if ats.get("readable_format"):
            score += 3
        if ats.get("no_tables_or_graphics"):
            score += 3
        return min(score, 15)

    def calculate_achievements_score(self, count: int) -> float:
        if count >= 5:
            return 10.0
        if count >= 3:
            return 7.0
        if count >= 1:
            return 4.0
        return 0.0

    def calculate_language_score(self, language: Dict[str, Any]) -> float:
        score = 0.0
        if language.get("grammar_accuracy", 0) >= 95:
            score += 4
        if language.get("clear_sentences"):
            score += 3
        if language.get("professional_tone"):
            score += 3
        return min(score, 10)

    def calculate_role_alignment_score(self, required: list, matched: list) -> float:
        if not required:
            return 0.0
        return min((len(matched) / len(required)) * 10, 10)

    def calculate_final_score(self, data: Dict[str, Any]) -> Dict[str, float]:
        final_score = (
            self.calculate_skills_score(
                data.get("required_skills") or [], data.get("candidate_skills") or []
            )
            + self.calculate_experience_score(
                float(data.get("required_experience") or 0),
                float(data.get("candidate_experience") or 0),
            )
            + self.calculate_ats_score(data.get("ats") or {})
            + self.calculate_achievements_score(int(data.get("achievement_count") or 0))
            + self.calculate_language_score(data.get("language") or {})
            + self.calculate_role_alignment_score(
                (data.get("role") or {}).get("required_keywords") or [],
                (data.get("role") or {}).get("matched_keywords") or [],
            )
        )
        return {"final_score": round(final_score, 2)}
