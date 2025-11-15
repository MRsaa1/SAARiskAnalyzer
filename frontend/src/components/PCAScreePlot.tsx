import { XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Line, ComposedChart, Bar } from 'recharts'

interface PCAScreePlotProps {
  data: Array<{ component: string; variance: number; cumulative: number }>
}

export default function PCAScreePlot({ data }: PCAScreePlotProps) {
  return (
    <ResponsiveContainer width="100%" height={300}>
      <ComposedChart data={data}>
        <CartesianGrid strokeDasharray="3 3" stroke="var(--border)" />
        <XAxis dataKey="component" stroke="var(--muted)" />
        <YAxis yAxisId="left" stroke="var(--muted)" />
        <YAxis yAxisId="right" orientation="right" stroke="var(--muted)" />
        <Tooltip 
          contentStyle={{ 
            background: 'var(--card)', 
            border: '1px solid var(--border)',
            borderRadius: '8px'
          }} 
        />
        <Bar yAxisId="left" dataKey="variance" fill="var(--accent)" name="Explained Variance" />
        <Line 
          yAxisId="right" 
          type="monotone" 
          dataKey="cumulative" 
          stroke="#22c55e" 
          strokeWidth={2}
          name="Cumulative %"
        />
      </ComposedChart>
    </ResponsiveContainer>
  )
}
