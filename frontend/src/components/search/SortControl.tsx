/**
 * SortControl Component
 * Flexible sorting control with separate sort_by and sort_order
 *
 * Features:
 * - Dropdown for sort field (price, departure_time, rating, popularity)
 * - Toggle button for sort direction (asc/desc)
 * - Visual icons for direction
 * - Aligned with backend refactored sorting
 */

import { ArrowUpDown, ArrowUp, ArrowDown } from 'lucide-react'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Button } from '@/components/ui/button'
import type { SearchSortBy, SearchSortOrder } from '@/types/search'

interface SortControlProps {
  sortBy: SearchSortBy
  sortOrder: SearchSortOrder
  onSortByChange: (value: SearchSortBy) => void
  onSortOrderChange: (value: SearchSortOrder) => void
}

const SORT_OPTIONS: Array<{ value: SearchSortBy; label: string; description: string }> = [
  {
    value: 'departure_time',
    label: 'Fecha de salida',
    description: 'Ordena por fecha y hora de salida',
  },
  {
    value: 'price',
    label: 'Precio',
    description: 'Ordena por precio por asiento',
  },
  {
    value: 'rating',
    label: 'Calificación',
    description: 'Ordena por calificación del conductor',
  },
  {
    value: 'popularity',
    label: 'Popularidad',
    description: 'Ordena por viajes más populares',
  },
]

export function SortControl({
  sortBy,
  sortOrder,
  onSortByChange,
  onSortOrderChange,
}: SortControlProps) {
  const toggleSortOrder = () => {
    onSortOrderChange(sortOrder === 'asc' ? 'desc' : 'asc')
  }

  const selectedOption = SORT_OPTIONS.find((opt) => opt.value === sortBy)

  return (
    <div className="space-y-2">
      <Label className="flex items-center gap-2">
        <ArrowUpDown className="h-4 w-4 text-gray-400" />
        Ordenar resultados
      </Label>

      <div className="flex gap-2">
        {/* Sort By Dropdown */}
        <Select value={sortBy} onValueChange={(value) => onSortByChange(value as SearchSortBy)}>
          <SelectTrigger className="flex-1">
            <SelectValue placeholder="Seleccionar orden" />
          </SelectTrigger>
          <SelectContent>
            {SORT_OPTIONS.map((option) => (
              <SelectItem key={option.value} value={option.value}>
                <div>
                  <div className="font-medium">{option.label}</div>
                  <div className="text-xs text-gray-500">{option.description}</div>
                </div>
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        {/* Sort Order Toggle */}
        <Button
          type="button"
          variant="outline"
          size="default"
          onClick={toggleSortOrder}
          className="px-4"
          title={sortOrder === 'asc' ? 'Orden ascendente' : 'Orden descendente'}
        >
          {sortOrder === 'asc' ? (
            <>
              <ArrowUp className="h-4 w-4 mr-2" />
              Asc
            </>
          ) : (
            <>
              <ArrowDown className="h-4 w-4 mr-2" />
              Desc
            </>
          )}
        </Button>
      </div>

      {/* Helper text */}
      {selectedOption && (
        <p className="text-xs text-gray-600">
          {sortOrder === 'asc' ? (
            <>
              {sortBy === 'price' && 'Menor a mayor precio'}
              {sortBy === 'departure_time' && 'Más pronto a más tarde'}
              {sortBy === 'rating' && 'Menor a mayor calificación'}
              {sortBy === 'popularity' && 'Menos a más popular'}
            </>
          ) : (
            <>
              {sortBy === 'price' && 'Mayor a menor precio'}
              {sortBy === 'departure_time' && 'Más tarde a más pronto'}
              {sortBy === 'rating' && 'Mayor a menor calificación'}
              {sortBy === 'popularity' && 'Más a menos popular'}
            </>
          )}
        </p>
      )}
    </div>
  )
}