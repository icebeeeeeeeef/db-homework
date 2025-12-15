package model.redpacket;

import java.util.HashMap;
import java.util.Map;

/**
 * 红包统计类，用于跟踪不同形状和大小的红包统计信息
 */
public class RedPacketStatistics {
    private final Map<String, Double> shapeAmounts = new HashMap<>();
    private final Map<String, Double> sizeAmounts = new HashMap<>();
    private final Map<String, Integer> shapeCounts = new HashMap<>();
    private final Map<String, Integer> sizeCounts = new HashMap<>();
    private final Map<String, Double> shapeSizeAmounts = new HashMap<>();
    private final Map<String, Integer> shapeSizeCounts = new HashMap<>();

    /**
     * 记录收集到的红包
     */
    public void recordCollected(RedPacket redPacket) {
        String shapeKey = redPacket.shape.getName();
        String sizeKey = redPacket.size.getName();
        String shapeSizeKey = shapeKey + sizeKey;

        // 更新形状统计
        shapeAmounts.merge(shapeKey, redPacket.amount, Double::sum);
        shapeCounts.merge(shapeKey, 1, Integer::sum);

        // 更新大小统计
        sizeAmounts.merge(sizeKey, redPacket.amount, Double::sum);
        sizeCounts.merge(sizeKey, 1, Integer::sum);

        // 更新形状+大小组合统计
        shapeSizeAmounts.merge(shapeSizeKey, redPacket.amount, Double::sum);
        shapeSizeCounts.merge(shapeSizeKey, 1, Integer::sum);
    }

    
    public Map<String, Double> getShapeAmounts() {
        return new HashMap<>(shapeAmounts);
    }

    public Map<String, Double> getSizeAmounts() {
        return new HashMap<>(sizeAmounts);
    }
    public Map<String, Integer> getShapeCounts() {
        return new HashMap<>(shapeCounts);
    }
    public Map<String, Integer> getSizeCounts() {
        return new HashMap<>(sizeCounts);
    }

    public Map<String, Double> getShapeSizeAmounts() {
        return new HashMap<>(shapeSizeAmounts);
    }

    public Map<String, Integer> getShapeSizeCounts() {
        return new HashMap<>(shapeSizeCounts);
    }

    public double getTotalAmount() {
        return shapeAmounts.values().stream().mapToDouble(Double::doubleValue).sum();
    }

    public int getTotalCount() {
        return shapeCounts.values().stream().mapToInt(Integer::intValue).sum();
    }
}
