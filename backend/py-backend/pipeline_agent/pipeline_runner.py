from google.adk.agents import Agent
from google.adk.models.lite_llm import LiteLlm
from pipeline_agent.rule_based_scoring import RuleBasedScoringEngine

# Dummy model for rule-based agent
dummy_model = LiteLlm(model="ollama_chat/llama3:8b")

class PipelineRunner:

    def __init__(self):
        self.scoring_engine = RuleBasedScoringEngine()

    def run(self, input_data):
        parsed_resume = input_data['parsed_resume']
        parsed_jd = input_data['parsed_jd']
        ats_result = input_data['ats_result']
        nlp_result = input_data['nlp_result']
        scoring_input = {
            "required_skills": parsed_jd["required_skills"],
            "candidate_skills": parsed_resume["skills"],

            "required_experience": parsed_jd["experience"],
            "candidate_experience": parsed_resume["experience"],

            "ats": ats_result,

            "achievement_count": parsed_resume["achievement_count"],

            "language": nlp_result,

            "role": {
                "required_keywords": parsed_jd["role_keywords"],
                "matched_keywords": parsed_resume["matched_role_keywords"]
            }
        }

        return self.scoring_engine.calculate_final_score(scoring_input)

class ScoringAgent(Agent):
    def __init__(self):
        super().__init__(
            name="RuleBasedScoringAgent",
            model=dummy_model,  # Dummy model since it's rule-based
            description="Rule-based ATS scoring agent",
            instruction="Calculate the final ATS score based on inputs.",
            tools=[]
        )

    def run(self, input_data):
        # Override run to use rule-based logic instead of LLM
        pipeline_runner = PipelineRunner()
        return pipeline_runner.run(input_data)

root_agent = ScoringAgent()
