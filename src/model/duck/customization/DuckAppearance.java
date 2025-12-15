package model.duck.customization;

import java.awt.Graphics2D;
import java.awt.Rectangle;

/**
 * Base contract for rendering a duck with optional accessories.
 */
public interface DuckAppearance {
    /**
     * Paints the duck inside the provided bounds.
     */
    void paint(Graphics2D g2, Rectangle bounds);
}
