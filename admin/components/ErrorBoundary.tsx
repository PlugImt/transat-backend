"use client";

import React from 'react';
import { AlertTriangle, RefreshCw, Copy, ChevronDown, ChevronRight } from 'lucide-react';

interface Props {
  children: React.ReactNode;
  fallback?: React.ReactNode;
  onError?: (error: Error, errorInfo: React.ErrorInfo) => void;
  level?: 'page' | 'component' | 'section';
}

interface State {
  hasError: boolean;
  error?: Error;
  errorInfo?: React.ErrorInfo;
  showDetails: boolean;
}

class ErrorBoundary extends React.Component<Props, State> {
  private retryCount = 0;
  private maxRetries = 3;

  constructor(props: Props) {
    super(props);
    this.state = { 
      hasError: false,
      showDetails: false
    };
  }

  static getDerivedStateFromError(error: Error): Partial<State> {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    this.setState({ errorInfo });
    
    // Log error details
    console.error('Error caught by boundary:', error, errorInfo);
    
    // Call custom error handler if provided
    this.props.onError?.(error, errorInfo);
    
    // In production, you could send this to an error reporting service
    if (process.env.NODE_ENV === 'production') {
      // Example: Send to error reporting service
      // errorReportingService.captureException(error, { extra: errorInfo });
    }
  }

  handleRetry = () => {
    if (this.retryCount < this.maxRetries) {
      this.retryCount++;
      this.setState({ 
        hasError: false, 
        error: undefined, 
        errorInfo: undefined,
        showDetails: false
      });
    } else {
      window.location.reload();
    }
  };

  handleCopyError = async () => {
    const errorDetails = `
Error: ${this.state.error?.message}
Stack: ${this.state.error?.stack}
Component Stack: ${this.state.errorInfo?.componentStack}
Timestamp: ${new Date().toISOString()}
User Agent: ${navigator.userAgent}
URL: ${window.location.href}
    `.trim();

    try {
      await navigator.clipboard.writeText(errorDetails);
      // Could show a toast notification here
    } catch (err) {
      console.error('Failed to copy error details:', err);
    }
  };

  toggleDetails = () => {
    this.setState(prev => ({ showDetails: !prev.showDetails }));
  };

  getErrorUI() {
    const { level = 'component' } = this.props;
    const { error, showDetails } = this.state;
    
    const isPageLevel = level === 'page';
    const containerClass = isPageLevel 
      ? "min-h-screen flex items-center justify-center p-8 bg-gray-50" 
      : "min-h-[200px] flex items-center justify-center p-8";
    
    return (
      <div className={containerClass}>
        <div className="text-center max-w-md mx-auto">
          <AlertTriangle className="h-12 w-12 text-red-500 mx-auto mb-4" />
          <h3 className="text-lg font-semibold text-gray-900 mb-2">
            {isPageLevel ? 'Erreur de chargement de la page' : 'Une erreur est survenue'}
          </h3>
          <p className="text-gray-600 mb-6">
            {error?.message || 'Quelque chose s&apos;est mal passé.'}
          </p>
          
          <div className="flex flex-col sm:flex-row gap-3 justify-center mb-4">
            <button
              onClick={this.handleRetry}
              className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
            >
              <RefreshCw className="h-4 w-4 mr-2" />
              {this.retryCount < this.maxRetries ? 'Réessayer' : 'Actualiser la page'}
            </button>
            
            <button
              onClick={this.handleCopyError}
              className="inline-flex items-center px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700 transition-colors"
            >
              <Copy className="h-4 w-4 mr-2" />
              Copier l&apos;erreur
            </button>
          </div>

          {process.env.NODE_ENV === 'development' && (
            <div className="text-left bg-gray-100 rounded-lg p-3">
              <button
                onClick={this.toggleDetails}
                className="flex items-center text-sm text-gray-700 hover:text-gray-900 mb-2"
              >
                {showDetails ? <ChevronDown className="h-4 w-4 mr-1" /> : <ChevronRight className="h-4 w-4 mr-1" />}
                Détails techniques
              </button>
              
              {showDetails && (
                <div className="text-xs text-gray-600 font-mono bg-white p-3 rounded border overflow-auto max-h-40">
                  <div className="mb-2">
                    <strong>Erreur:</strong> {error?.message}
                  </div>
                  {error?.stack && (
                    <div className="mb-2">
                      <strong>Stack:</strong>
                      <pre className="whitespace-pre-wrap mt-1">{error.stack}</pre>
                    </div>
                  )}
                  {this.state.errorInfo?.componentStack && (
                    <div>
                      <strong>Component Stack:</strong>
                      <pre className="whitespace-pre-wrap mt-1">{this.state.errorInfo.componentStack}</pre>
                    </div>
                  )}
                </div>
              )}
            </div>
          )}
          
          <p className="text-xs text-gray-500 mt-4">
            Tentatives: {this.retryCount}/{this.maxRetries}
          </p>
        </div>
      </div>
    );
  }

  render() {
    if (this.state.hasError) {
      return this.props.fallback || this.getErrorUI();
    }

    return this.props.children;
  }
}

export default ErrorBoundary;