import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import '@testing-library/jest-dom';
import SurveyForm from './SurveyForm';
import { api } from '../services/api';

// Mock the API service
jest.mock('../services/api', () => ({
  api: {
    get: jest.fn(),
    post: jest.fn()
  }
}));

// Mock API responses
const mockSurvey = {
  id: 'test-survey-123',
  title: 'Customer Satisfaction Survey',
  description: 'Please provide your feedback',
  questions: [
    {
      id: 'q1',
      text: 'How would you rate our service?',
      type: 'RATING',
      required: true,
      options: {
        min: 1,
        max: 5
      }
    },
    {
      id: 'q2',
      text: 'What did you like most about our service?',
      type: 'TEXT',
      required: false
    },
    {
      id: 'q3',
      text: 'Which of our products have you used?',
      type: 'MULTIPLE_CHOICE',
      required: true,
      options: {
        choices: ['Product A', 'Product B', 'Product C'],
        multiple: true
      }
    }
  ]
};

const renderWithRouter = (component) => {
  return render(
    <BrowserRouter>
      {component}
    </BrowserRouter>
  );
};

describe('SurveyForm Component', () => {
  beforeEach(() => {
    // Reset mock API functions
    api.get.mockReset();
    api.post.mockReset();
    
    // Setup API mock response
    api.get.mockResolvedValue({ data: mockSurvey });
    api.post.mockResolvedValue({ data: { id: 'response-123' } });
  });

  test('renders survey form with all questions', async () => {
    renderWithRouter(<SurveyForm surveyId="test-survey-123" />);
    
    // Wait for data to load
    await waitFor(() => {
      expect(screen.getByText('Customer Satisfaction Survey')).toBeInTheDocument();
    });
    
    // Check if the API was called with the correct survey ID
    expect(api.get).toHaveBeenCalledWith('/surveys/test-survey-123');
    
    // Check if all questions are rendered
    expect(screen.getByText('How would you rate our service?')).toBeInTheDocument();
    expect(screen.getByText('What did you like most about our service?')).toBeInTheDocument();
    expect(screen.getByText('Which of our products have you used?')).toBeInTheDocument();
    
    // Check required indicator
    expect(screen.getAllByText('*').length).toBe(2); // Two required questions
  });

  test('displays validation error when required fields are not filled', async () => {
    renderWithRouter(<SurveyForm surveyId="test-survey-123" />);
    
    // Wait for data to load
    await waitFor(() => {
      expect(screen.getByText('Customer Satisfaction Survey')).toBeInTheDocument();
    });
    
    // Submit form without filling required fields
    const submitButton = screen.getByText('Submit');
    fireEvent.click(submitButton);
    
    // Check for validation error messages
    await waitFor(() => {
      expect(screen.getByText('Please rate our service')).toBeInTheDocument();
      expect(screen.getByText('Please select at least one product')).toBeInTheDocument();
    });
    
    // API should not be called
    expect(api.post).not.toHaveBeenCalled();
  });

  test('successfully submits form with valid data', async () => {
    renderWithRouter(<SurveyForm surveyId="test-survey-123" />);
    
    // Wait for data to load
    await waitFor(() => {
      expect(screen.getByText('Customer Satisfaction Survey')).toBeInTheDocument();
    });
    
    // Fill in the form
    // Rating question - select 4
    const ratingInput = screen.getByLabelText('4');
    fireEvent.click(ratingInput);
    
    // Text question
    const textInput = screen.getByPlaceholderText('Your answer');
    fireEvent.change(textInput, { target: { value: 'Great service!' } });
    
    // Multiple choice question - select Product A and Product C
    const checkboxA = screen.getByLabelText('Product A');
    const checkboxC = screen.getByLabelText('Product C');
    fireEvent.click(checkboxA);
    fireEvent.click(checkboxC);
    
    // Submit the form
    const submitButton = screen.getByText('Submit');
    fireEvent.click(submitButton);
    
    // Check if API was called with the correct data
    await waitFor(() => {
      expect(api.post).toHaveBeenCalledWith('/responses', {
        surveyId: 'test-survey-123',
        answers: [
          { questionId: 'q1', value: '4' },
          { questionId: 'q2', value: 'Great service!' },
          { questionId: 'q3', value: JSON.stringify(['Product A', 'Product C']) }
        ]
      });
    });
    
    // Check success message
    expect(screen.getByText('Thank you for your response!')).toBeInTheDocument();
  });

  test('handles API error when loading survey', async () => {
    // Mock API error
    api.get.mockRejectedValue(new Error('Failed to load survey'));
    
    renderWithRouter(<SurveyForm surveyId="test-survey-123" />);
    
    // Check for error message
    await waitFor(() => {
      expect(screen.getByText('Error loading survey: Failed to load survey')).toBeInTheDocument();
    });
  });

  test('handles API error when submitting response', async () => {
    // Setup API to succeed for get but fail for post
    api.get.mockResolvedValue({ data: mockSurvey });
    api.post.mockRejectedValue(new Error('Failed to submit response'));
    
    renderWithRouter(<SurveyForm surveyId="test-survey-123" />);
    
    // Wait for data to load
    await waitFor(() => {
      expect(screen.getByText('Customer Satisfaction Survey')).toBeInTheDocument();
    });
    
    // Fill in the form with valid data
    const ratingInput = screen.getByLabelText('4');
    fireEvent.click(ratingInput);
    
    const checkboxA = screen.getByLabelText('Product A');
    fireEvent.click(checkboxA);
    
    // Submit the form
    const submitButton = screen.getByText('Submit');
    fireEvent.click(submitButton);
    
    // Check for error message
    await waitFor(() => {
      expect(screen.getByText('Error submitting response: Failed to submit response')).toBeInTheDocument();
    });
  });
}); 