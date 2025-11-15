import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts'

interface RollingVaRChartProps {
  data: Array<{ date: string; var: number; cvar: number }>
}

export default function RollingVaRChart({ data }: RollingVaRChartProps) {
  return (
    <ResponsiveContainer width="100%" height={300}>
      <LineChart data={data}>
        <CartesianGrid strokeDasharray="3 3" stroke="var(--border)" />
        <XAxis dataKey="date" stroke="var(--muted)" />
        <YAxis stroke="var(--muted)" />
        <Tooltip 
          contentStyle={{ 
            background: 'var(--card)', 
            border: '1px solid var(--border)',
            borderRadius: '8px',
            color: 'var(--fg)'
          }} 
        />
        <Legend />
        <Line 
          type="monotone" 
          dataKey="var" 
          stroke="var(--accent)" 
          strokeWidth={2}
          name="VaR (99%)"
          dot={false}
        />
        <Line 
          type="monotone" 
          dataKey="cvar" 
          stroke="#f97316" 
          strokeWidth={2}
          name="CVaR (99%)"
          dot={false}
        />
      </LineChart>
    </ResponsiveContainer>
  )
}
