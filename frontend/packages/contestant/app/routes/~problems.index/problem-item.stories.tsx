import type { Meta, StoryObj } from "@storybook/react";
import { ProblemItem } from "./problem-item";

export default {
  title: "pages/problems/ProblemItem",
} satisfies Meta;

type Story = StoryObj;

export const Default: Story = {
  render: () => (
    <div className="grid grid-cols-2 gap-64">
      <ProblemItem
        code="ABC"
        title="あいしーてぃーえすしーだよあああああああ"
        score={{
          maxScore: 200,
          score: 200,
          rawScore: 200,
          penalty: 0,
          fullScore: true,
          rawFullScore: true,
        }}
      />
      <ProblemItem
        code="ABC"
        title="あいしーてぃーえすしーだよあああああああ"
        score={{
          maxScore: 200,
          score: 160,
          rawScore: 200,
          penalty: -40,
          fullScore: false,
          rawFullScore: true,
        }}
      />
      <ProblemItem
        code="ABC"
        title="あいしーてぃーえすしーだよあああああああ"
        score={{
          maxScore: 200,
          score: 100,
          rawScore: 120,
          penalty: -20,
          fullScore: false,
          rawFullScore: false,
        }}
      />
      <ProblemItem
        code="ABC"
        title="あいしーてぃーえすしーだよあああああああ"
        score={{
          maxScore: 200,
        }}
      />
    </div>
  ),
};
