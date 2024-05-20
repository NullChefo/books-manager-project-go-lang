## Abstraction for Working with Databases

Develop a package for generating SQL queries (Query builder) that includes an implementation of the public API described in the table below. For greater flexibility and convenience, the arguments and return types are deliberately omitted.

<table>
    <thead>
        <tr>
            <th width="220">Name</th>
            <th width="480">Description</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>
                <sub>createTable()</sub>
            </td>
            <td>
                Creates a new table based on the provided structure
            </td>
        </tr>
        <tr>
            <td>
                <sub>select()</sub>
            </td>
            <td>
                Creates a new SELECT query based on the structure and returns an array of the same structure
            </td>
        </tr>
        <tr>
            <td>
                <sub>insert()</sub>
            </td>
            <td>
                Creates a new INSERT query based on the structure and returns an object
            </td>
        </tr>
        <tr>
            <td>
                <sub>update()</sub>
            </td>
            <td>
                Creates a new UPDATE query based on the structure and returns an object
            </td>
        </tr>
    </tbody>
</table>

## WEB API for Book Management

Develop a WEB API for managing books. A book consists of the following attributes:
- title
- ISBN
- author
- year of publication

Implement the following public API endpoints.

<table>
    <thead>
        <tr>
            <th width="220">Name</th>
            <th width="480">Description</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>
                <strong>GET</strong>
                <sub>/books</sub>
            </td>
            <td>
                Returns a collection of all books recorded in the system. (Consider a variant where the system returns a limited number of books)
            </td>
        </tr>
        <tr>
            <td>
                <strong>GET</strong>
                <sub>/books</sub>
            </td>
            <td>
                Returns data for a single book
            </td>
        </tr>
        <tr>
            <td>
                <strong>POST</strong>
                <sub>/books</sub>
            </td>
            <td>
                Creates a new book
            </td>
        </tr>
        <tr>
            <td>
                <strong>PUT</strong>
                <sub>/books/:id</sub>
            </td>
            <td>
                Updates the information of an already added book in the database
            </td>
        </tr>
        <tr>
            <td>
                <strong>DELETE</strong>
                <sub>/books/:id</sub>
            </td>
            <td>
                Deletes a book by identifier
            </td>
        </tr>
    </tbody>
</table>

### Project Submission

After developing the project, send a link to your GitHub repository in the assignment posted in Google Classroom. The deadline for the project is the end of the semester.