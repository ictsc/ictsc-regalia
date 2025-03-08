import type { Meta, StoryObj } from "@storybook/react";
import { ProblemsPage } from "./page";

export default {
  title: "pages/problems",
  component: ProblemsPage,
} satisfies Meta<typeof ProblemsPage>;

type Story = StoryObj<typeof ProblemsPage>;

export const Default: Story = {
  args: {
    problems: [
      ...Array.from({ length: 4 }, () => ({
        $typeName: "contestant.v1.Problem" as const,
        code: "ABC",
        title: "問題 ABC",
        maxScore: 200,
        category: "Network",
      })),
      ...Array.from({ length: 5 }, () => ({
        $typeName: "contestant.v1.Problem" as const,
        code: "ABC",
        title: "問題 ABC",
        maxScore: 200,
        category: "Server",
        score: {
          $typeName: "contestant.v1.Score" as const,
          score: 100,
          markedScore: 120,
          penalty: -20,
        },
      })),
      ...Array.from({ length: 2 }, () => ({
        $typeName: "contestant.v1.Problem" as const,
        code: "ABC",
        title: "問題 ABC",
        maxScore: 200,
        category: "Network",
        score: {
          $typeName: "contestant.v1.Score" as const,
          score: 160,
          markedScore: 200,
          penalty: -40,
        },
      })),
      ...Array.from({ length: 5 }, () => ({
        $typeName: "contestant.v1.Problem" as const,
        code: "ABC",
        title: "問題 ABC",
        maxScore: 200,
        category: "Server",
        score: {
          $typeName: "contestant.v1.Score" as const,
          score: 200,
          markedScore: 200,
          penalty: 0,
        },
      })),
    ],
  },
};
