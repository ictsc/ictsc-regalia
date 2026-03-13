import type { Meta, StoryObj } from "@storybook/react";
import { create } from "@bufbuild/protobuf";
import { timestampFromDate } from "@bufbuild/protobuf/wkt";
import { ScheduleEntrySchema } from "@ictsc/proto/contestant/v1";
import { IndexPage } from "./index.page";

export default {
  title: "pages/index",
  component: IndexPage,
} satisfies Meta<typeof IndexPage>;

type Story = StoryObj<typeof IndexPage>;

const day1Am = create(ScheduleEntrySchema, {
  name: "day1-am",
  startAt: timestampFromDate(new Date("2026-01-01T10:00:00+09:00")),
  endAt: timestampFromDate(new Date("2026-01-01T12:00:00+09:00")),
});

const day1Pm = create(ScheduleEntrySchema, {
  name: "day1-pm",
  startAt: timestampFromDate(new Date("2026-01-01T13:00:00+09:00")),
  endAt: timestampFromDate(new Date("2099-12-31T23:59:59+09:00")),
});

const day2Am = create(ScheduleEntrySchema, {
  name: "day2-am",
  startAt: timestampFromDate(new Date("2100-01-01T10:00:00+09:00")),
  endAt: timestampFromDate(new Date("2100-01-01T12:00:00+09:00")),
});

const entries = [day1Am, day1Pm, day2Am];

export const InContest: Story = {
  name: "競技中",
  args: {
    state: "in_contest",
    currentScheduleName: "day1-pm",
    nextScheduleName: "day2-am",
    entries,
  },
};

export const Waiting: Story = {
  name: "競技時間外（待機中）",
  args: {
    state: "waiting",
    nextScheduleName: "day1-am",
    entries,
  },
};

export const Ended: Story = {
  name: "競技終了",
  args: {
    state: "ended",
    entries,
  },
};
