"use client";

import { ChevronDown, X } from "lucide-react";
import { useEffect, useRef, useState } from "react";

interface TagAutocompleteProps {
  options: { value: string; label: string }[];
  selectedTags: string[];
  onChange: (tags: string[]) => void;
  placeholder?: string;
  className?: string;
  id?: string;
}

export default function TagAutocomplete({
  options,
  selectedTags,
  onChange,
  placeholder = "Rechercher et sélectionner...",
  className = "",
  id,
}: TagAutocompleteProps) {
  const [inputValue, setInputValue] = useState("");
  const [isOpen, setIsOpen] = useState(false);
  const [filteredOptions, setFilteredOptions] = useState(options);
  const inputRef = useRef<HTMLInputElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const filtered = options.filter(
      (option) =>
        option.label.toLowerCase().includes(inputValue.toLowerCase()) &&
        !selectedTags.includes(option.value),
    );
    setFilteredOptions(filtered);
  }, [inputValue, options, selectedTags]);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
        setIsOpen(false);
        setInputValue("");
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const addTag = (value: string) => {
    if (!selectedTags.includes(value)) {
      onChange([...selectedTags, value]);
    }
    setInputValue("");
    setIsOpen(false);
  };

  const removeTag = (value: string) => {
    onChange(selectedTags.filter((tag) => tag !== value));
  };

  const handleInputKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Backspace" && inputValue === "" && selectedTags.length > 0) {
      removeTag(selectedTags[selectedTags.length - 1]);
    }
    if (e.key === "Enter") {
      e.preventDefault();
      if (filteredOptions.length > 0) {
        addTag(filteredOptions[0].value);
      }
    }
    if (e.key === "Escape") {
      setIsOpen(false);
      setInputValue("");
    }
  };

  return (
    <div ref={containerRef} className={`relative ${className}`}>
      <div
        className={`w-full min-h-[40px] px-2 sm:px-3 py-2 border border-gray-300 rounded-md focus-within:outline-none focus-within:ring-2 focus-within:ring-blue-500 focus-within:border-blue-500 bg-white cursor-text ${
          isOpen ? "ring-2 ring-blue-500 border-blue-500" : ""
        }`}
        role="combobox"
        tabIndex={0}
        aria-expanded={isOpen}
        onClick={() => {
          setIsOpen(true);
          inputRef.current?.focus();
        }}
        onKeyDown={(e) => {
          if (e.key === "Enter" || e.key === " ") {
            e.preventDefault();
            setIsOpen(true);
            inputRef.current?.focus();
          }
        }}
      >
        <div className="flex flex-wrap items-center gap-1">
          {selectedTags.map((tag) => {
            const option = options.find((opt) => opt.value === tag);
            return (
              <span
                key={tag}
                className="inline-flex items-center px-2 sm:px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800"
              >
                {option?.label || tag}
                <button
                  type="button"
                  onClick={(e) => {
                    e.stopPropagation();
                    removeTag(tag);
                  }}
                  className="ml-1 hover:text-blue-600"
                >
                  <X className="h-3 w-3" />
                </button>
              </span>
            );
          })}
          <div className="flex-1 flex items-center min-w-[100px]">
            <input
              ref={inputRef}
              id={id}
              type="text"
              value={inputValue}
              onChange={(e) => setInputValue(e.target.value)}
              onKeyDown={handleInputKeyDown}
              onFocus={() => setIsOpen(true)}
              placeholder={selectedTags.length === 0 ? placeholder : ""}
              className="flex-1 border-none outline-none bg-transparent text-sm"
            />
            <ChevronDown
              className={`h-4 w-4 text-gray-400 transition-transform ${isOpen ? "rotate-180" : ""}`}
            />
          </div>
        </div>
      </div>

      {isOpen && (
        <div className="absolute z-50 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg max-h-60 overflow-auto">
          {filteredOptions.length === 0 ? (
            <div className="px-3 py-2 text-sm text-gray-500">
              {inputValue ? "Aucun résultat trouvé" : "Tous les rôles sont déjà sélectionnés"}
            </div>
          ) : (
            filteredOptions.map((option) => (
              <button
                key={option.value}
                type="button"
                onClick={() => addTag(option.value)}
                className="w-full px-3 py-2 text-left text-sm hover:bg-gray-100 focus:bg-gray-100 focus:outline-none"
              >
                {option.label}
              </button>
            ))
          )}
        </div>
      )}
    </div>
  );
}
