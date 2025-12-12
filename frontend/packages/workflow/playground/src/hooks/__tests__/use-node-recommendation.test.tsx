/**
 * Copyright 2025 Coze Studio. All rights reserved.
 *
 * useNodeRecommendation Hook 测试
 */

import { renderHook, act, waitFor } from '@testing-library/react';
import { useNodeRecommendation } from '../use-node-recommendation';

// Mock RecommendationService
jest.mock('../../services/recommendation-service', () => {
  return {
    RecommendationService: jest.fn().mockImplementation(() => ({
      getRecommendations: jest.fn().mockResolvedValue({
        recommendations: [
          {
            nodeType: 'LoopNode',
            score: 0.95,
            reason: 'Array output requires iteration',
            metadata: { confidence: 'high' },
          },
          {
            nodeType: 'DatabaseQueryNode',
            score: 0.88,
            reason: 'Common data processing pattern',
            metadata: { confidence: 'medium' },
          },
        ],
        metadata: {
          requestId: 'test-request-123',
          totalCandidates: 10,
          timestamp: new Date().toISOString(),
        },
      }),
      recordFeedback: jest.fn().mockResolvedValue({ success: true }),
    })),
  };
});

describe('useNodeRecommendation', () => {
  const defaultOptions = {
    workflowId: 'test-workflow-id',
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should initialize with empty recommendations', () => {
    const { result } = renderHook(() => useNodeRecommendation(defaultOptions));

    expect(result.current.recommendations).toEqual([]);
    expect(result.current.loading).toBe(false);
    expect(result.current.error).toBeNull();
    expect(result.current.requestId).toBe('');
  });

  it('should load recommendations successfully', async () => {
    const { result } = renderHook(() => useNodeRecommendation(defaultOptions));

    await act(async () => {
      await result.current.loadRecommendations({
        sourceNodeId: 'node-1',
        sourceNodeType: 'CodeNode',
        limit: 5,
      });
    });

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.recommendations).toHaveLength(2);
    expect(result.current.recommendations[0].nodeType).toBe('LoopNode');
    expect(result.current.recommendations[0].score).toBe(0.95);
    expect(result.current.requestId).toBe('test-request-123');
    expect(result.current.error).toBeNull();
  });

  it('should set loading state during fetch', async () => {
    const { result } = renderHook(() => useNodeRecommendation(defaultOptions));

    let loadingStates: boolean[] = [];

    act(() => {
      result.current.loadRecommendations({
        sourceNodeId: 'node-1',
        sourceNodeType: 'CodeNode',
      });
      loadingStates.push(result.current.loading);
    });

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    loadingStates.push(result.current.loading);
    expect(loadingStates).toContain(true);
    expect(result.current.loading).toBe(false);
  });

  it('should handle errors gracefully', async () => {
    const mockError = new Error('Network error');
    
    // Mock service to throw error
    const mockService = require('../../services/recommendation-service').RecommendationService;
    mockService.mockImplementationOnce(() => ({
      getRecommendations: jest.fn().mockRejectedValue(mockError),
    }));

    const { result } = renderHook(() => useNodeRecommendation(defaultOptions));

    await act(async () => {
      await result.current.loadRecommendations({
        sourceNodeId: 'node-1',
        sourceNodeType: 'CodeNode',
      });
    });

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.error).toEqual(mockError);
    expect(result.current.recommendations).toEqual([]);
  });

  it('should submit feedback successfully', async () => {
    const { result } = renderHook(() => useNodeRecommendation(defaultOptions));

    // First load recommendations to get requestId
    await act(async () => {
      await result.current.loadRecommendations({
        sourceNodeId: 'node-1',
        sourceNodeType: 'CodeNode',
      });
    });

    await waitFor(() => {
      expect(result.current.requestId).toBe('test-request-123');
    });

    // Then submit feedback
    await act(async () => {
      await result.current.submitFeedback('LoopNode', 'selected');
    });

    const mockService = require('../../services/recommendation-service').RecommendationService;
    const serviceInstance = mockService.mock.results[0].value;
    
    expect(serviceInstance.recordFeedback).toHaveBeenCalledWith({
      requestId: 'test-request-123',
      selectedNode: 'LoopNode',
      userAction: 'selected',
    });
  });

  it('should clear recommendations', async () => {
    const { result } = renderHook(() => useNodeRecommendation(defaultOptions));

    // Load recommendations
    await act(async () => {
      await result.current.loadRecommendations({
        sourceNodeId: 'node-1',
        sourceNodeType: 'CodeNode',
      });
    });

    await waitFor(() => {
      expect(result.current.recommendations).toHaveLength(2);
    });

    // Clear recommendations
    act(() => {
      result.current.clearRecommendations();
    });

    expect(result.current.recommendations).toEqual([]);
    expect(result.current.requestId).toBe('');
    expect(result.current.error).toBeNull();
  });

  it('should reload with last parameters', async () => {
    const { result } = renderHook(() => useNodeRecommendation(defaultOptions));

    const params = {
      sourceNodeId: 'node-1',
      sourceNodeType: 'CodeNode',
      limit: 5,
    };

    // Initial load
    await act(async () => {
      await result.current.loadRecommendations(params);
    });

    await waitFor(() => {
      expect(result.current.recommendations).toHaveLength(2);
    });

    // Clear recommendations
    act(() => {
      result.current.clearRecommendations();
    });

    expect(result.current.recommendations).toEqual([]);

    // Reload
    await act(async () => {
      await result.current.reload();
    });

    await waitFor(() => {
      expect(result.current.recommendations).toHaveLength(2);
    });

    const mockService = require('../../services/recommendation-service').RecommendationService;
    const serviceInstance = mockService.mock.results[0].value;
    
    // Should be called twice with same parameters
    expect(serviceInstance.getRecommendations).toHaveBeenCalledTimes(2);
  });

  it('should not submit feedback without requestId', async () => {
    const { result } = renderHook(() => useNodeRecommendation(defaultOptions));

    // Try to submit feedback without loading recommendations first
    await act(async () => {
      await result.current.submitFeedback('LoopNode', 'selected');
    });

    const mockService = require('../../services/recommendation-service').RecommendationService;
    const serviceInstance = mockService.mock.results[0].value;
    
    expect(serviceInstance.recordFeedback).not.toHaveBeenCalled();
  });

  it('should handle all feedback actions', async () => {
    const { result } = renderHook(() => useNodeRecommendation(defaultOptions));

    // Load recommendations
    await act(async () => {
      await result.current.loadRecommendations({
        sourceNodeId: 'node-1',
        sourceNodeType: 'CodeNode',
      });
    });

    await waitFor(() => {
      expect(result.current.requestId).toBeTruthy();
    });

    const mockService = require('../../services/recommendation-service').RecommendationService;
    const serviceInstance = mockService.mock.results[0].value;

    // Test all user actions
    await act(async () => {
      await result.current.submitFeedback('LoopNode', 'selected');
    });
    expect(serviceInstance.recordFeedback).toHaveBeenLastCalledWith(
      expect.objectContaining({ userAction: 'selected' })
    );

    await act(async () => {
      await result.current.submitFeedback('LoopNode', 'dismissed');
    });
    expect(serviceInstance.recordFeedback).toHaveBeenLastCalledWith(
      expect.objectContaining({ userAction: 'dismissed' })
    );

    await act(async () => {
      await result.current.submitFeedback('LoopNode', 'ignored');
    });
    expect(serviceInstance.recordFeedback).toHaveBeenLastCalledWith(
      expect.objectContaining({ userAction: 'ignored' })
    );
  });
});
