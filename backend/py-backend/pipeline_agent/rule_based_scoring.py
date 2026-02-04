from typing import Dict

class RuleBasedScoringEngine:

    def calculate_skills_score(self, required, candidate):
        if not required:
            return 0
        matched = len(set(required).intersection(set(candidate)))
        return min((matched / len(required)) * 30, 30)

    def calculate_experience_score(self, required_years, candidate_years):
        if required_years <= 0:
            return 0
        return min((candidate_years / required_years) * 25, 25)

    def calculate_ats_score(self, ats):
        score = 0
        if ats.get("keyword_density_ok"):
            score += 5
        if ats.get("standard_headings"):
            score += 4
        if ats.get("readable_format"):
            score += 3
        if ats.get("no_tables_or_graphics"):
            score += 3
        return min(score, 15)

    def calculate_achievements_score(self, count):
        if count >= 5:
            return 10
        elif count >= 3:
            return 7
        elif count >= 1:
            return 4
        return 0

    def calculate_language_score(self, language):
        score = 0
        if language.get("grammar_accuracy", 0) >= 95:
            score += 4
        if language.get("clear_sentences"):
            score += 3
        if language.get("professional_tone"):
            score += 3
        return min(score, 10)

    def calculate_role_alignment_score(self, required, matched):
        if not required:
            return 0
        return min((len(matched) / len(required)) * 10, 10)

    def calculate_final_score(self, data):
        final_score = (
            self.calculate_skills_score(
                data["required_skills"], data["candidate_skills"]
            )
            + self.calculate_experience_score(
                data["required_experience"], data["candidate_experience"]
            )
            + self.calculate_ats_score(data["ats"])
            + self.calculate_achievements_score(data["achievement_count"])
            + self.calculate_language_score(data["language"])
            + self.calculate_role_alignment_score(
                data["role"]["required_keywords"],
                data["role"]["matched_keywords"]
            )
        )

        return {
            "final_score": round(final_score, 2)
        }


# ✅ ADK TOOL FUNCTION (THIS IS WHAT ADK CALLS)
def rule_based_scoring_tool(scoring_input: Dict) -> Dict:
    engine = RuleBasedScoringEngine()
    return engine.calculate_final_score(scoring_input)
