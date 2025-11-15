import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, ReferenceLine } from 'recharts'

interface PnLHistogramProps {
  data: Array<{ bucket: string; count: number }>
  var: number
  cvar: number
}

export default function PnLHistogram({ data, var: varValue, cvar }: PnLHistogramProps) {
  return (
    <ResponsiveContainer width="100%" height={300}>
      <BarChart data={data}>
        <CartesianGrid strokeDasharray="3 3" stroke="var(--border)" />
        <XAxis dataKey="bucket" stroke="var(--muted)" />
        <YAxis stroke="var(--muted)" />
        <Tooltip 
          contentStyle={{ 
            background: 'var(--card)', 
            border: '1px solid var(--border)',
            borderRadius: '8px' 
          }} 
        />
        <Bar dataKey="count" fill="var(--accent)" />
        <ReferenceLine x={varValue} stroke="#ef4444" strokeWidth={2} label="VaR" />
        <ReferenceLine x={cvar} stroke="#dc2626" strokeWidth={2} strokeDasharray="5 5" label="CVaR" />
      </BarChart>
    </ResponsiveContainer>
  )
}
