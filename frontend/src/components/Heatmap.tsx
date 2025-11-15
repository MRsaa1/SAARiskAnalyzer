interface HeatmapProps {
  data: number[][]
  labels: string[]
}

export default function Heatmap({ data, labels }: HeatmapProps) {
  const cellSize = 60
  
  const getColor = (value: number) => {
    if (value > 0.7) return '#22c55e'
    if (value > 0.3) return '#84cc16'
    if (value > -0.3) return '#fbbf24'
    if (value > -0.7) return '#f97316'
    return '#ef4444'
  }
  
  return (
    <div className="overflow-x-auto">
      <svg 
        width={(labels.length + 1) * cellSize} 
        height={(labels.length + 1) * cellSize}
      >
        {labels.map((label, i) => (
          <text
            key={`top-${i}`}
            x={(i + 1) * cellSize + cellSize / 2}
            y={cellSize - 10}
            textAnchor="middle"
            fill="var(--fg)"
            fontSize="12"
            fontWeight="600"
          >
            {label}
          </text>
        ))}
        
        {labels.map((label, i) => (
          <text
            key={`left-${i}`}
            x={cellSize - 10}
            y={(i + 1) * cellSize + cellSize / 2}
            textAnchor="end"
            dominantBaseline="middle"
            fill="var(--fg)"
            fontSize="12"
            fontWeight="600"
          >
            {label}
          </text>
        ))}
        
        {data.map((row, i) =>
          row.map((value, j) => (
            <g key={`cell-${i}-${j}`}>
              <rect
                x={(j + 1) * cellSize}
                y={(i + 1) * cellSize}
                width={cellSize}
                height={cellSize}
                fill={getColor(value)}
                stroke="var(--border)"
                strokeWidth="1"
                opacity="0.8"
              />
              <text
                x={(j + 1) * cellSize + cellSize / 2}
                y={(i + 1) * cellSize + cellSize / 2}
                textAnchor="middle"
                dominantBaseline="middle"
                fill="#000"
                fontSize="11"
                fontWeight="600"
              >
                {value.toFixed(2)}
              </text>
            </g>
          ))
        )}
      </svg>
    </div>
  )
}
