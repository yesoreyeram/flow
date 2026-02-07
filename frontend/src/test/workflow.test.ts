import { describe, it, expect } from 'vitest';
import { httpRequestSchema, transformSchema, conditionSchema } from '../types/workflow';

describe('Workflow Type Validation', () => {
  describe('httpRequestSchema', () => {
    it('validates valid HTTP request config', () => {
      const validConfig = {
        url: 'https://api.example.com/data',
        method: 'GET' as const,
        headers: { 'Content-Type': 'application/json' },
      };

      const result = httpRequestSchema.safeParse(validConfig);
      expect(result.success).toBe(true);
    });

    it('rejects invalid URL', () => {
      const invalidConfig = {
        url: 'not-a-url',
        method: 'GET' as const,
      };

      const result = httpRequestSchema.safeParse(invalidConfig);
      expect(result.success).toBe(false);
    });
  });

  describe('transformSchema', () => {
    it('validates valid transform config', () => {
      const validConfig = {
        code: 'return input.map(x => x * 2)',
        language: 'javascript' as const,
      };

      const result = transformSchema.safeParse(validConfig);
      expect(result.success).toBe(true);
    });
  });

  describe('conditionSchema', () => {
    it('validates valid condition config', () => {
      const validConfig = {
        conditions: [
          {
            field: 'status',
            operator: 'equals' as const,
            value: 'success',
          },
        ],
        combinator: 'AND' as const,
      };

      const result = conditionSchema.safeParse(validConfig);
      expect(result.success).toBe(true);
    });
  });
});
