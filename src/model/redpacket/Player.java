package model.redpacket;

import geom.Vector2;

public class Player {
    public Vector2 pos;
    public final double radius;

    public Player(double x, double y, double radius) {
        this.pos = new Vector2(x, y);
        this.radius = radius;
    }

    public boolean collide(RedPacket rp) {
        return !rp.collected && pos.distance(rp.pos) <= radius;
    }
}

