package attendance.model;

public class Student {
    private Integer id;
    private String studentNo;
    private String name;
    private String photoPath;

    public Student(Integer id, String studentNo, String name, String photoPath) {
        this.id = id;
        this.studentNo = studentNo;
        this.name = name;
        this.photoPath = photoPath;
    }

    public Student(String studentNo, String name, String photoPath) {
        this(null, studentNo, name, photoPath);
    }

    public Integer getId() {
        return id;
    }

    public void setId(Integer id) {
        this.id = id;
    }

    public String getStudentNo() {
        return studentNo;
    }

    public String getName() {
        return name;
    }

    public String getPhotoPath() {
        return photoPath;
    }
}
