package model.duck.customization;

import java.awt.Graphics2D;
import java.awt.Rectangle;

/**
 * Helper base class for duck accessories.
 */
public abstract class AccessoryDecorator implements DuckAppearance {
    protected final DuckAppearance inner;

    protected AccessoryDecorator(DuckAppearance inner) {
        this.inner = inner;
    }

    @Override
    public void paint(Graphics2D g2, Rectangle bounds) {
        inner.paint(g2, bounds);
        paintAccessory(g2, bounds);
    }

    protected abstract void paintAccessory(Graphics2D g2, Rectangle bounds);
}
