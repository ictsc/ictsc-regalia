type JsonRecord = Record<string, unknown>;

export function toCamelCase<T>(value: unknown): T {
  return mapKeys(value, snakeToCamel) as T;
}

export function toSnakeCase<T>(value: unknown): T {
  return mapKeys(value, camelToSnake) as T;
}

function mapKeys(value: unknown, convert: (key: string) => string): unknown {
  if (Array.isArray(value)) return value.map((item) => mapKeys(item, convert));
  if (value == null || typeof value !== "object" || value instanceof Date) {
    return value;
  }
  return Object.fromEntries(
    Object.entries(value as JsonRecord).map(([key, item]) => [
      convert(key),
      mapKeys(item, convert),
    ]),
  );
}

function snakeToCamel(value: string): string {
  return value.replace(/_([a-z0-9])/g, (_, char: string) => char.toUpperCase());
}

function camelToSnake(value: string): string {
  return value.replace(/[A-Z]/g, (char) => `_${char.toLowerCase()}`);
}
