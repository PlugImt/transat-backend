"use client";

import { Check, ChevronDown, Search, X } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { useClickOutside } from "@/lib/hooks";

interface ComboboxOption {
  value: string;
  label: string;
}

interface ComboboxProps {
  options: ComboboxOption[];
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  label?: string;
  searchPlaceholder?: string;
  emptyText?: string;
  clearable?: boolean;
}

export default function Combobox({
  options,
  value,
  onChange,
  placeholder = "Sélectionner...",
  label,
  searchPlaceholder = "Rechercher...",
  emptyText = "Aucun résultat",
  clearable = true,
}: ComboboxProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState("");
  const searchInputRef = useRef<HTMLInputElement>(null);

  // Generate unique ID for accessibility
  const buttonId = `combobox-${Math.random().toString(36).substr(2, 9)}`;

  const containerRef = useClickOutside<HTMLDivElement>(() => setIsOpen(false));

  const filteredOptions = options.filter((option) =>
    option.label.toLowerCase().includes(searchTerm.toLowerCase()),
  );

  const selectedOption = options.find((option) => option.value === value);

  const handleSelect = (optionValue: string) => {
    onChange(optionValue);
    setIsOpen(false);
    setSearchTerm("");
  };

  const handleClear = (e: React.MouseEvent) => {
    e.stopPropagation();
    onChange("");
    setSearchTerm("");
  };

  const handleToggle = () => {
    setIsOpen(!isOpen);
    if (!isOpen) {
      setTimeout(() => {
        searchInputRef.current?.focus();
      }, 0);
    }
  };

  useEffect(() => {
    if (isOpen) {
      setSearchTerm("");
    }
  }, [isOpen]);

  return (
    <div className="relative" ref={containerRef}>
      {label && (
        <label htmlFor={buttonId} className="block text-sm font-medium text-gray-700 mb-2">
          {label}
        </label>
      )}

      <div className="relative">
        <div className="relative flex items-center">
          <button
            id={buttonId}
            type="button"
            onClick={handleToggle}
            className="w-full flex items-center justify-between px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          >
            <span
              className={`block truncate ${!selectedOption ? "text-gray-500" : "text-gray-900"}`}
            >
              {selectedOption ? selectedOption.label : placeholder}
            </span>
            <div className="flex items-center space-x-1">
              <ChevronDown
                className={`h-4 w-4 text-gray-400 transition-transform ${isOpen ? "rotate-180" : ""}`}
              />
            </div>
          </button>
          {clearable && selectedOption && (
            <button
              type="button"
              onClick={handleClear}
              className="absolute right-8 p-1 hover:bg-gray-100 rounded transition-colors z-10"
              title="Effacer la sélection"
            >
              <X className="h-3 w-3 text-gray-400 hover:text-gray-600" />
            </button>
          )}
        </div>

        {isOpen && (
          <div className="absolute z-50 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg">
            {/* Search input */}
            <div className="p-2 border-b border-gray-200">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                <input
                  ref={searchInputRef}
                  type="text"
                  placeholder={searchPlaceholder}
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
              </div>
            </div>

            {/* Options list */}
            <div className="max-h-60 overflow-auto">
              {filteredOptions.length === 0 ? (
                <div className="px-3 py-2 text-sm text-gray-500 text-center">{emptyText}</div>
              ) : (
                filteredOptions.map((option) => (
                  <button
                    key={option.value}
                    type="button"
                    onClick={() => handleSelect(option.value)}
                    className={`w-full flex items-center justify-between px-3 py-2 text-left hover:bg-gray-50 transition-colors ${
                      option.value === value ? "bg-blue-50 text-blue-900" : "text-gray-900"
                    }`}
                  >
                    <span className="block truncate">{option.label}</span>
                    {option.value === value && <Check className="h-4 w-4 text-blue-600" />}
                  </button>
                ))
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
