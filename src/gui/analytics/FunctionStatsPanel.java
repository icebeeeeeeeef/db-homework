package gui;

import javax.swing.*;
import javax.swing.table.DefaultTableModel;
import java.awt.*;
import java.util.List;

public class FunctionStatsPanel extends JPanel {
    private final DefaultTableModel model;

    public FunctionStatsPanel(List<CodeStatsResult.FunctionStat> stats) {
        setLayout(new BorderLayout());
        String[] columns = {"Language", "Functions", "Avg", "Min", "Max", "Median"};
        model = new DefaultTableModel(columns, 0) {
            @Override
            public boolean isCellEditable(int row, int column) {
                return false;
            }
        };
        JTable table = new JTable(model);
        table.setRowHeight(28);
        table.setFont(table.getFont().deriveFont(14f));
        table.getTableHeader().setFont(table.getFont().deriveFont(Font.BOLD, 14f));
        add(new JScrollPane(table), BorderLayout.CENTER);
        refresh(stats);
    }

    public void refresh(List<CodeStatsResult.FunctionStat> stats) {
        model.setRowCount(0);
        if (stats == null) return;
        for (CodeStatsResult.FunctionStat stat : stats) {
            model.addRow(new Object[]{
                    stat.language,
                    stat.count,
                    String.format("%.2f", stat.average),
                    stat.min,
                    stat.max,
                    String.format("%.2f", stat.median)
            });
        }
    }
}
