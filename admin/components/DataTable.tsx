"use client";

import {
  type ColumnDef,
  type ColumnFiltersState,
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  type SortingState,
  useReactTable,
} from "@tanstack/react-table";
import clsx from "clsx";
import { ChevronDown, ChevronUp, Filter, Search } from "lucide-react";
import { memo, useEffect, useMemo, useState } from "react";
import { useDebounce } from "@/lib/hooks";

interface DataTableProps<T> {
  data: T[];
  columns: ColumnDef<T>[];
  searchPlaceholder?: string;
  globalFilterColumn?: string;
  className?: string;
  onRowClick?: (row: T) => void;
  // Filter options
  showFilter?: boolean;
  filterActive?: boolean;
  onFilterToggle?: () => void;
  filterLabel?: string;
  filterActiveLabel?: string;
}

function DataTable<T>({
  data,
  columns,
  searchPlaceholder = "Rechercher...",
  globalFilterColumn: _globalFilterColumn,
  className,
  onRowClick,
  // Filter props
  showFilter = false,
  filterActive = false,
  onFilterToggle,
  filterLabel = "Filtre",
  filterActiveLabel = "Filtre actif",
}: DataTableProps<T>) {
  const [sorting, setSorting] = useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);
  const [globalFilter, setGlobalFilter] = useState("");
  const [searchInput, setSearchInput] = useState("");

  // Débounce la recherche pour améliorer les performances
  const debouncedSearchTerm = useDebounce(searchInput, 300);

  // Mettre à jour le filtre global quand le terme de recherche débounce change
  useEffect(() => {
    setGlobalFilter(debouncedSearchTerm);
  }, [debouncedSearchTerm]);

  // Memoize table configuration to avoid recreating on every render
  const tableConfig = useMemo(
    () => ({
      data,
      columns,
      getCoreRowModel: getCoreRowModel(),
      getSortedRowModel: getSortedRowModel(),
      getFilteredRowModel: getFilteredRowModel(),
      getPaginationRowModel: getPaginationRowModel(),
      onSortingChange: setSorting,
      onColumnFiltersChange: setColumnFilters,
      onGlobalFilterChange: setGlobalFilter,
      // Utilise la fonction de filtrage par défaut pour rechercher dans toutes les colonnes
      state: {
        sorting,
        columnFilters,
        globalFilter,
      },
      initialState: {
        pagination: {
          pageSize: 10,
        },
      },
    }),
    [data, columns, sorting, columnFilters, globalFilter],
  );

  const table = useReactTable(tableConfig);

  return (
    <div className={clsx("space-y-4", className)}>
      {/* Search and Filters */}
      <div className="flex flex-col sm:flex-row sm:items-center gap-3 mb-2">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
          <input
            placeholder={searchPlaceholder}
            value={searchInput}
            onChange={(e) => setSearchInput(e.target.value)}
            className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          />
          {/* Indicateur de recherche en cours */}
          {searchInput !== debouncedSearchTerm && (
            <div className="absolute right-3 top-1/2 transform -translate-y-1/2">
              <div className="animate-spin rounded-full h-4 w-4 border-2 border-blue-600 border-t-transparent"></div>
            </div>
          )}
        </div>

        {/* Filter button */}
        {showFilter && onFilterToggle && (
          <button
            type="button"
            onClick={onFilterToggle}
            className={`inline-flex items-center px-3 py-2 rounded-md text-sm font-medium transition-colors border ${
              filterActive
                ? "bg-orange-100 text-orange-800 border-orange-200 hover:bg-orange-50"
                : "bg-white text-gray-700 border-gray-300 hover:bg-gray-50"
            }`}
            title={filterActive ? filterActiveLabel : filterLabel}
          >
            <Filter className="h-4 w-4 mr-2" />
            <span className="hidden sm:inline">
              {filterActive ? filterActiveLabel : filterLabel}
            </span>
          </button>
        )}
      </div>

      {/* Active filter indicator */}
      {showFilter && filterActive && (
        <div className="mb-3">
          <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-orange-100 text-orange-700 border border-orange-200">
            <Filter className="h-3 w-3 mr-1" />
            {filterActiveLabel}
          </span>
        </div>
      )}

      {/* Table */}
      <div className="overflow-hidden border border-gray-200 rounded-lg">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              {table.getHeaderGroups().map((headerGroup) => (
                <tr key={headerGroup.id}>
                  {headerGroup.headers.map((header) => (
                    <th
                      key={header.id}
                      className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                    >
                      {header.isPlaceholder ? null : (
                        <div
                          {...{
                            className: header.column.getCanSort()
                              ? "cursor-pointer select-none flex items-center space-x-1"
                              : "",
                            onClick: header.column.getToggleSortingHandler(),
                          }}
                        >
                          <span>
                            {flexRender(header.column.columnDef.header, header.getContext())}
                          </span>
                          {header.column.getCanSort() && (
                            <span className="flex flex-col">
                              {header.column.getIsSorted() === "asc" ? (
                                <ChevronUp className="h-4 w-4" />
                              ) : header.column.getIsSorted() === "desc" ? (
                                <ChevronDown className="h-4 w-4" />
                              ) : (
                                <div className="h-4 w-4 opacity-50">
                                  <ChevronUp className="h-2 w-4" />
                                  <ChevronDown className="h-2 w-4" />
                                </div>
                              )}
                            </span>
                          )}
                        </div>
                      )}
                    </th>
                  ))}
                </tr>
              ))}
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {table.getRowModel().rows.map((row) => (
                <tr
                  key={row.id}
                  className={clsx("hover:bg-gray-50", onRowClick && "cursor-pointer")}
                  onClick={() => onRowClick?.(row.original)}
                >
                  {row.getVisibleCells().map((cell) => (
                    <td key={cell.id} className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Pagination */}
      <div className="flex items-center justify-between">
        <div className="text-sm text-gray-700">
          {(() => {
            const pageIndex = table.getState().pagination.pageIndex;
            const pageSize = table.getState().pagination.pageSize;
            const totalRows = table.getFilteredRowModel().rows.length;
            const start = totalRows === 0 ? 0 : pageIndex * pageSize + 1;
            const end = Math.min((pageIndex + 1) * pageSize, totalRows);

            return `Affichage de ${start} à ${end} sur ${totalRows} résultats`;
          })()}
        </div>
        <div className="flex items-center space-x-2">
          <button
            type="button"
            className="px-3 py-2 text-sm bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
            onClick={() => table.previousPage()}
            disabled={!table.getCanPreviousPage()}
          >
            Précédent
          </button>
          <span className="text-sm text-gray-700">
            Page {table.getState().pagination.pageIndex + 1} sur {table.getPageCount()}
          </span>
          <button
            type="button"
            className="px-3 py-2 text-sm bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
            onClick={() => table.nextPage()}
            disabled={!table.getCanNextPage()}
          >
            Suivant
          </button>
        </div>
      </div>
    </div>
  );
}

// Export memoized component with proper typing
export default memo(DataTable) as <T = unknown>(props: DataTableProps<T>) => React.ReactElement;
