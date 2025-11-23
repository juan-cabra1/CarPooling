/**
 * RadiusSlider Component
 * Allows user to select search radius for geospatial queries
 *
 * Features:
 * - Range: 1 km - 100 km
 * - Default: 10 km
 * - Shows current value
 * - Can be disabled when no location is selected
 */

import { useState, useEffect } from 'react'
import { Maximize2 } from 'lucide-react'
import { Label } from '@/components/ui/label'
import { Slider } from '@/components/ui/slider'

interface RadiusSliderProps {
  value: number // in kilometers
  onChange: (radius: number) => void
  disabled?: boolean
  min?: number // default: 1
  max?: number // default: 100
  label?: string
}

export function RadiusSlider({
  value,
  onChange,
  disabled = false,
  min = 1,
  max = 100,
  label = 'Radio de búsqueda',
}: RadiusSliderProps) {
  const [localValue, setLocalValue] = useState(value)

  // Sync with prop value
  useEffect(() => {
    setLocalValue(value)
  }, [value])

  const handleChange = (newValue: number[]) => {
    const radius = newValue[0]
    setLocalValue(radius)
    onChange(radius)
  }

  return (
    <div className="space-y-2">
      <div className="flex items-center justify-between">
        <Label className="flex items-center gap-2">
          <Maximize2 className="h-4 w-4 text-gray-400" />
          {label}
        </Label>
        <span className="text-sm font-medium text-gray-700">
          {localValue} km
        </span>
      </div>

      <Slider
        value={[localValue]}
        onValueChange={handleChange}
        min={min}
        max={max}
        step={1}
        disabled={disabled}
        className="w-full"
      />

      <div className="flex justify-between text-xs text-gray-500">
        <span>{min} km</span>
        <span>{max} km</span>
      </div>

      {disabled && (
        <p className="text-xs text-gray-500 italic">
          Selecciona una ubicación para habilitar el radio de búsqueda
        </p>
      )}
    </div>
  )
}