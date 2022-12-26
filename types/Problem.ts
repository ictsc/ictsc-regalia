export interface Problem {
    id: string;
    code: string;
    title: string;
    body: string;
    point: number;
    solved_criterion: number | null;
    previous_problem_id: string | null;
    author_id: string;
}